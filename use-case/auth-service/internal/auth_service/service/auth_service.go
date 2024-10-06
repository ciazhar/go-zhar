package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ciazhar/go-zhar/use-case/auth-service/internal/auth_service/model"
	"github.com/ciazhar/go-zhar/use-case/auth-service/internal/auth_service/repository"
	"github.com/ciazhar/go-zhar/use-case/auth-service/pkg/bcrypt"
	"github.com/ciazhar/go-zhar/use-case/auth-service/pkg/jwt"
	"github.com/ciazhar/go-zhar/use-case/auth-service/pkg/logger"
)

type AuthService struct {
	userPGRepo    *repository.UsersPostgresRepository
	authRedisRepo *repository.AuthRedisRepository
}

func NewAuthService(userRepo *repository.UsersPostgresRepository, authRedisRepo *repository.AuthRedisRepository) *AuthService {
	return &AuthService{
		userPGRepo:    userRepo,
		authRedisRepo: authRedisRepo,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, user model.User) error {

	// Hash the password
	hashedPassword, err := bcrypt.HashPassword(user.Password)
	if err != nil {
		return errors.New(fmt.Sprintf("could not hash password: %v", err))
	}
	user.Password = hashedPassword

	tx, err := s.userPGRepo.BeginTransaction(ctx)
	if err != nil {
		return errors.New("could not begin transaction")
	}
	defer tx.Rollback(ctx)

	// Proceed with inserting the user and other operations
	err = s.userPGRepo.Insert(ctx, user)
	if err != nil {
		return errors.New("could not insert user")
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.New("could not commit transaction")
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, body model.LoginRequest) (model.LoginResponse, error) {

	// Retrieve user from the database
	user, err := s.userPGRepo.GetByUsername(ctx, body.Username)
	if err != nil {
		return model.LoginResponse{}, logger.LogAndReturnError(ctx, err,
			"could not get user", map[string]string{
				"username": body.Username,
			})
	}

	// Check password
	if !bcrypt.CheckPasswordHash(body.Password, user.Password) {
		return model.LoginResponse{}, logger.LogAndReturnWarning(ctx, nil,
			"invalid credentials", map[string]string{
				"username": body.Username,
			})
	}

	// Generate tokens
	accessToken, err := jwt.GenerateJWTToken(user.ID, jwt.AccessTokenTTL)
	if err != nil {
		return model.LoginResponse{}, logger.LogAndReturnError(ctx, err,
			"could not generate access token",
			map[string]string{
				"userID": user.ID,
			})
	}

	refreshToken, err := jwt.GenerateJWTToken(user.ID, jwt.RefreshTokenTTL)
	if err != nil {
		return model.LoginResponse{}, logger.LogAndReturnError(ctx, err,
			"could not generate refresh token",
			map[string]string{
				"userID": user.ID,
			})
	}

	// Store JWT in Redis (allow multiple tokens per user)
	err = s.authRedisRepo.StoreAccessToken(ctx, user.ID, accessToken)
	if err != nil {
		return model.LoginResponse{}, logger.LogAndReturnError(ctx, err,
			"could not store access token",
			map[string]string{
				"userID":      user.ID,
				"accessToken": accessToken,
			})
	}

	// Store refresh token in Redis with expiration
	err = s.authRedisRepo.StoreRefreshToken(ctx, user.ID, refreshToken)
	if err != nil {
		return model.LoginResponse{}, logger.LogAndReturnError(ctx, err,
			"could not store refresh token",
			map[string]string{
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
		return model.RefreshTokenResponse{}, errors.New(fmt.Sprintf("could not validate refresh token: %v", err))
	}

	// Check if refresh token exists in Redis
	if err = s.authRedisRepo.IsRefreshTokenExist(ctx, claims.UserID, refreshToken); err != nil {
		return model.RefreshTokenResponse{}, errors.New(fmt.Sprintf("could not validate refresh token: %v", err))
	}

	// Generate new access token
	newAccessToken, err := jwt.GenerateJWTToken(claims.UserID, jwt.AccessTokenTTL)
	if err != nil {
		return model.RefreshTokenResponse{}, errors.New(fmt.Sprintf("could not generate new access token: %v", err))
	}

	// Store JWT in Redis (allow multiple tokens per user)
	if err = s.authRedisRepo.StoreAccessToken(ctx, claims.UserID, newAccessToken); err != nil {
		return model.RefreshTokenResponse{}, errors.New(fmt.Sprintf("could not store access token: %v", err))
	}

	return model.RefreshTokenResponse{
		AccessToken: newAccessToken,
	}, nil
}

func (s *AuthService) Protected(ctx context.Context, accessToken string) (string, error) {

	// Validate token
	claims, err := jwt.ValidateJWT(accessToken)
	if err != nil {
		return "", errors.New(fmt.Sprintf("could not validate access token: %v", err))
	}

	// Check if the token exists in Redis
	exists, err := s.authRedisRepo.IsAccessTokenExist(ctx, claims.UserID, accessToken)
	if err != nil {
		return "", errors.New(fmt.Sprintf("could not validate access token: %v", err))
	}
	if !exists {
		return "", errors.New("access token not found")
	}

	return claims.UserID, nil
}

func (s *AuthService) Logout(ctx context.Context, accessToken string) error {

	// Validate token
	claims, err := jwt.ValidateJWT(accessToken)
	if err != nil {
		return errors.New(fmt.Sprintf("could not validate access token: %v", err))
	}

	// Remove the specific access token from Redis
	err = s.authRedisRepo.RemoveAccessToken(ctx, claims.UserID, accessToken)
	if err != nil {
		return errors.New(fmt.Sprintf("could not remove access token: %v", err))
	}

	// Delete refresh token from Redis
	err = s.authRedisRepo.RemoveRefreshToken(ctx, claims.UserID)
	if err != nil {
		return errors.New(fmt.Sprintf("could not remove refresh token: %v", err))
	}

	return nil
}

func (s *AuthService) Revoke(ctx context.Context, accessToken string) error {

	// Validate token
	claims, err := jwt.ValidateJWT(accessToken)
	if err != nil {
		return errors.New(fmt.Sprintf("could not validate access token: %v", err))
	}

	// Remove all access tokens for the user
	err = s.authRedisRepo.RemoveAllAccessTokens(ctx, claims.UserID)
	if err != nil {
		return errors.New(fmt.Sprintf("could not revoke tokens: %v", err))
	}

	// Delete all refresh tokens for the user
	err = s.authRedisRepo.RemoveRefreshToken(ctx, claims.UserID)
	if err != nil {
		return errors.New(fmt.Sprintf("could not revoke tokens: %v", err))
	}

	return nil
}
