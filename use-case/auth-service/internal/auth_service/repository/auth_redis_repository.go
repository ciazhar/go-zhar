package repository

import (
	"context"
	"errors"
	"github.com/ciazhar/go-zhar/use-case/auth-service/pkg/jwt"
	"github.com/go-redis/redis/v8"
)

const (
	accessTokenPrefix  = "auth:access_token:"
	refreshTokenPrefix = "auth:refresh_token:"
)

type AuthRedisRepository struct {
	redis *redis.Client
}

func NewAuthRedisRepository(redis *redis.Client) *AuthRedisRepository {
	return &AuthRedisRepository{redis: redis}
}

func (r *AuthRedisRepository) StoreAccessToken(ctx context.Context, userId string, accessToken string) error {

	err := r.redis.SAdd(ctx, accessTokenPrefix+userId, accessToken).Err()
	if err != nil {
		return err
	}

	err = r.redis.Expire(ctx, accessTokenPrefix+userId, jwt.RefreshTokenTTL).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRedisRepository) StoreRefreshToken(ctx context.Context, userId string, refreshToken string) error {
	return r.redis.Set(ctx, refreshTokenPrefix+userId, refreshToken, jwt.RefreshTokenTTL).Err()
}

func (r *AuthRedisRepository) GetAccessToken(ctx context.Context, userId string) (string, error) {
	return r.redis.SRandMember(ctx, accessTokenPrefix+userId).Result()
}

func (r *AuthRedisRepository) GetRefreshToken(ctx context.Context, userId string) (string, error) {
	return r.redis.Get(ctx, refreshTokenPrefix+userId).Result()
}

func (r *AuthRedisRepository) IsRefreshTokenExist(ctx context.Context, userId string, refreshToken string) error {
	val, err := r.redis.Get(ctx, refreshTokenPrefix+userId).Result()
	if err != nil {
		return err
	}

	if val != refreshToken {
		return errors.New("invalid refresh token")
	}

	return nil
}

func (r *AuthRedisRepository) IsAccessTokenExist(ctx context.Context, userId string, accessToken string) (bool, error) {
	exists, err := r.redis.SIsMember(ctx, accessTokenPrefix+userId, accessToken).Result()
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *AuthRedisRepository) RemoveAccessToken(ctx context.Context, userId string, accessToken string) error {
	return r.redis.SRem(ctx, accessTokenPrefix+userId, accessToken).Err()
}

func (r *AuthRedisRepository) RemoveRefreshToken(ctx context.Context, userId string) error {
	return r.redis.Del(ctx, refreshTokenPrefix+userId).Err()
}

func (r *AuthRedisRepository) RemoveAllAccessTokens(ctx context.Context, userId string) error {
	return r.redis.Del(ctx, accessTokenPrefix+userId).Err()
}
