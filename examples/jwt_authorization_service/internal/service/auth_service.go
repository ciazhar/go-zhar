package service

import (
	"context"

	"github.com/ciazhar/go-start-small/examples/jwt_authorization_service/internal/model"
	"github.com/ciazhar/go-start-small/examples/jwt_authorization_service/internal/repository"
	"github.com/ciazhar/go-start-small/pkg/hashing/bcrypt"
	"github.com/ciazhar/go-start-small/pkg/jwt"
	"github.com/ciazhar/go-start-small/pkg/logger"
	timeutil "github.com/ciazhar/go-start-small/pkg/time_util"
	"github.com/spf13/viper"
)

type AuthServiceInterface interface {
	RegisterUser(ctx context.Context, user model.User) error
	Login(ctx context.Context, body model.LoginRequest) (model.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (model.RefreshTokenResponse, error)
	Protected(ctx context.Context, accessToken string) (string, error)
	Logout(ctx context.Context, accessToken string) error
	Revoke(ctx context.Context, accessToken string) error
}

type AuthService struct {
	userPGRepo    repository.UsersPostgresRepositoryInterface
	authRedisRepo repository.AuthRedisRepositoryInterface
}

func NewAuthService(userRepo repository.UsersPostgresRepositoryInterface, authRedisRepo repository.AuthRedisRepositoryInterface) *AuthService {
	return &AuthService{
		userPGRepo:    userRepo,
		authRedisRepo: authRedisRepo,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, user model.User) error {

	// Hash the password
	hashedPassword, err := bcrypt.HashPassword(user.Password)
	if err != nil {
		return logger.LogAndReturnError(ctx, err, "could not hash password", nil)
	}
	user.Password = hashedPassword

	tx, err := s.userPGRepo.BeginTransaction(ctx)
	if err != nil {
		return logger.LogAndReturnError(ctx, err, "could not begin transaction", nil)
	}
	defer tx.Rollback(ctx)

	// Proceed with inserting the user and other operations
	err = s.userPGRepo.Insert(ctx, tx, user)
	if err != nil {
		return logger.LogAndReturnError(ctx, err, "could not insert user", nil)
	}

	if err = tx.Commit(ctx); err != nil {
		return logger.LogAndReturnError(ctx, err, "could not commit transaction", nil)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, body model.LoginRequest) (model.LoginResponse, error) {

	// Retrieve user from the database
	user, err := s.userPGRepo.GetByUsername(ctx, body.Username)
	if err != nil {
		return model.LoginResponse{}, logger.LogAndReturnError(ctx, err,
			"could not get user", map[string]interface{}{
				"username": body.Username,
			})
	}

	// Check password
	if !bcrypt.ValidatePassword(body.Password, user.Password) {
		return model.LoginResponse{}, logger.LogAndReturnWarning(ctx, nil,
			"invalid credentials", map[string]interface{}{
				"username": body.Username,
			})
	}

	// Parse access token TTL
	duration, err := timeutil.ParseTimeDuration(viper.GetString("jwt.accessTokenTTL"))
	if err != nil {
		return model.LoginResponse{}, logger.LogAndReturnError(ctx, err, "could not validate access token", nil)
	}

	// Generate tokens
	accessToken, err := jwt.GenerateJWTToken(user.ID, duration)
	if err != nil {
		return model.LoginResponse{}, logger.LogAndReturnError(ctx, err,
			"could not generate access token",
			map[string]interface{}{
				"userID": user.ID,
			})
	}

	// Parse access token TTL
	duration, err = timeutil.ParseTimeDuration(viper.GetString("jwt.refreshTokenTTL"))
	if err != nil {
		return model.LoginResponse{}, logger.LogAndReturnError(ctx, err, "could not validate refresh token", nil)
	}

	refreshToken, err := jwt.GenerateJWTToken(user.ID, duration)
	if err != nil {
		return model.LoginResponse{}, logger.LogAndReturnError(ctx, err,
			"could not generate refresh token",
			map[string]interface{}{
				"userID": user.ID,
			})
	}

	// Store JWT in Redis (allow multiple tokens per user)
	err = s.authRedisRepo.StoreAccessToken(ctx, user.ID, accessToken)
	if err != nil {
		return model.LoginResponse{}, logger.LogAndReturnError(ctx, err,
			"could not store access token",
			map[string]interface{}{
				"userID":      user.ID,
				"accessToken": accessToken,
			})
	}

	// Store refresh token in Redis with expiration
	err = s.authRedisRepo.StoreRefreshToken(ctx, user.ID, refreshToken)
	if err != nil {
		return model.LoginResponse{}, logger.LogAndReturnError(ctx, err,
			"could not store refresh token",
			map[string]interface{}{
				"userID":       user.ID,
				"refreshToken": refreshToken,
			})
	}

	return model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (model.RefreshTokenResponse, error) {

	// Validate refresh token
	claims, err := jwt.ValidateJWT(refreshToken)
	if err != nil {
		return model.RefreshTokenResponse{}, logger.LogAndReturnError(ctx, err, "could not validate refresh token", nil)
	}

	// Check if refresh token exists in Redis
	if err = s.authRedisRepo.IsRefreshTokenExist(ctx, claims.UserID, refreshToken); err != nil {
		return model.RefreshTokenResponse{}, logger.LogAndReturnError(ctx, err, "could not validate refresh token", nil)
	}

	// Parse access token TTL
	duration, err := timeutil.ParseTimeDuration(viper.GetString("jwt.accessTokenTTL"))
	if err != nil {
		return model.RefreshTokenResponse{}, logger.LogAndReturnError(ctx, err, "error parsing access token TTL", nil)
	}

	// Generate new access token
	newAccessToken, err := jwt.GenerateJWTToken(claims.UserID, duration)
	if err != nil {
		return model.RefreshTokenResponse{}, logger.LogAndReturnError(ctx, err, "could not generate new access token", nil)
	}

	// Store JWT in Redis (allow multiple tokens per user)
	if err = s.authRedisRepo.StoreAccessToken(ctx, claims.UserID, newAccessToken); err != nil {
		return model.RefreshTokenResponse{}, logger.LogAndReturnError(ctx, err, "could not store access token", nil)
	}

	return model.RefreshTokenResponse{
		AccessToken: newAccessToken,
	}, nil
}

func (s *AuthService) Protected(ctx context.Context, accessToken string) (string, error) {

	// Validate token
	claims, err := jwt.ValidateJWT(accessToken)
	if err != nil {
		return "", logger.LogAndReturnError(ctx, err, "could not validate access token", nil)
	}

	// Check if the token exists in Redis
	exists, err := s.authRedisRepo.IsAccessTokenExist(ctx, claims.UserID, accessToken)
	if err != nil {
		return "", logger.LogAndReturnError(ctx, err, "could not validate access token", nil)
	}
	if !exists {
		return "", logger.LogAndReturnError(ctx, err, "access token not found", nil)
	}

	return claims.UserID, nil
}

func (s *AuthService) Logout(ctx context.Context, accessToken string) error {

	// Validate token
	claims, err := jwt.ValidateJWT(accessToken)
	if err != nil {
		return logger.LogAndReturnError(ctx, err, "could not validate access token", nil)
	}

	// Remove the specific access token from Redis
	err = s.authRedisRepo.RemoveAccessToken(ctx, claims.UserID, accessToken)
	if err != nil {
		return logger.LogAndReturnError(ctx, err, "could not remove access token", nil)
	}

	// Delete refresh token from Redis
	err = s.authRedisRepo.RemoveRefreshToken(ctx, claims.UserID)
	if err != nil {
		return logger.LogAndReturnError(ctx, err, "could not remove refresh token", nil)
	}

	return nil
}

func (s *AuthService) Revoke(ctx context.Context, accessToken string) error {
	// Validate token
	claims, err := jwt.ValidateJWT(accessToken)
	if err != nil {
		return logger.LogAndReturnError(ctx, err, "could not validate access token", nil)
	}

	// Remove all access tokens for the user
	err = s.authRedisRepo.RemoveAllAccessTokens(ctx, claims.UserID)
	if err != nil {
		return logger.LogAndReturnError(ctx, err, "could not revoke access tokens", nil)
	}

	// Delete refresh token for the user
	err = s.authRedisRepo.RemoveRefreshToken(ctx, claims.UserID)
	if err != nil {
		return logger.LogAndReturnError(ctx, err, "could not revoke refresh token", nil)
	}

	return nil
}
