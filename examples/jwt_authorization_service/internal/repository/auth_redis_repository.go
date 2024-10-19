package repository

import (
	"context"
	"errors"

	timeutil "github.com/ciazhar/go-start-small/pkg/time_util"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

const (
	accessTokenPrefix  = "auth:access_token:"
	refreshTokenPrefix = "auth:refresh_token:"
)

type AuthRedisRepositoryInterface interface {
	StoreAccessToken(ctx context.Context, userId string, accessToken string) error
	StoreRefreshToken(ctx context.Context, userId string, refreshToken string) error
	GetAccessToken(ctx context.Context, userId string) (string, error)
	GetRefreshToken(ctx context.Context, userId string) (string, error)
	IsAccessTokenExist(ctx context.Context, userId string, accessToken string) (bool, error)
	IsRefreshTokenExist(ctx context.Context, userId string, refreshToken string) error
	RemoveAccessToken(ctx context.Context, userId string, accessToken string) error
	RemoveRefreshToken(ctx context.Context, userId string) error
	RemoveAllAccessTokens(ctx context.Context, userId string) error
}

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

	duration, err := timeutil.ParseTimeDuration(viper.GetString("jwt.refreshTokenTTL"))
	if err != nil {
		return err
	}

	err = r.redis.Expire(ctx, accessTokenPrefix+userId, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRedisRepository) StoreRefreshToken(ctx context.Context, userId string, refreshToken string) error {

	duration, err := timeutil.ParseTimeDuration(viper.GetString("jwt.refreshTokenTTL"))
	if err != nil {
		return err
	}

	return r.redis.Set(ctx, refreshTokenPrefix+userId, refreshToken, duration).Err()
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
