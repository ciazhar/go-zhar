package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type HashUserProfileRepository struct {
	redis *redis.Client
}

func NewHashUserProfileRepository(redisClient *redis.Client) *HashUserProfileRepository {
	return &HashUserProfileRepository{
		redis: redisClient,
	}
}

// SetUserProfile Set user profile details
func (r *HashUserProfileRepository) SetUserProfile(ctx context.Context, userID string, profile map[string]interface{}) error {
	key := fmt.Sprintf("user:%s:profile", userID)
	return r.redis.HSet(ctx, key, profile).Err()
}

// GetUserProfile Get user profile details
func (r *HashUserProfileRepository) GetUserProfile(ctx context.Context, userID string) (map[string]string, error) {
	key := fmt.Sprintf("user:%s:profile", userID)
	return r.redis.HGetAll(ctx, key).Result()
}
