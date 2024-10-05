package repository

import (
	"context"
	"errors"
	"github.com/ciazhar/go-zhar/use-case/auth-service/pkg/jwt"
	"github.com/go-redis/redis/v8"
)

type AuthRedisRepository struct {
	redis *redis.Client
}

func NewAuthRedisRepository(redis *redis.Client) *AuthRedisRepository {
	return &AuthRedisRepository{redis: redis}
}

func (r *AuthRedisRepository) StoreAccessToken(ctx context.Context, userId string, accessToken string) error {
	return r.redis.SAdd(ctx, "user:"+userId, accessToken).Err()
}

func (r *AuthRedisRepository) StoreRefreshToken(ctx context.Context, userId string, refreshToken string) error {
	return r.redis.Set(ctx, "refresh:"+userId, refreshToken, jwt.RefreshTokenTTL).Err()
}

func (r *AuthRedisRepository) GetAccessToken(ctx context.Context, userId string) (string, error) {
	return r.redis.SRandMember(ctx, "user:"+userId).Result()
}

func (r *AuthRedisRepository) GetRefreshToken(ctx context.Context, userId string) (string, error) {
	return r.redis.Get(ctx, "refresh:"+userId).Result()
}

func (r *AuthRedisRepository) IsRefreshTokenExist(ctx context.Context, userId string, refreshToken string) error {
	val, err := r.redis.Get(ctx, "refresh:"+userId).Result()
	if err != nil {
		return err
	}

	if val != refreshToken {
		return errors.New("invalid refresh token")
	}

	return nil
}

func (r *AuthRedisRepository) IsAccessTokenExist(ctx context.Context, userId string, accessToken string) (bool, error) {
	exists, err := r.redis.SIsMember(ctx, "user:"+userId, accessToken).Result()
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *AuthRedisRepository) RemoveAccessToken(ctx context.Context, userId string, accessToken string) error {
	return r.redis.SRem(ctx, "user:"+userId, accessToken).Err()
}

func (r *AuthRedisRepository) RemoveRefreshToken(ctx context.Context, userId string) error {
	return r.redis.Del(ctx, "refresh:"+userId).Err()
}

func (r *AuthRedisRepository) RemoveAllAccessTokens(ctx context.Context, userId string) error {
	return r.redis.Del(ctx, "user:"+userId).Err()
}
