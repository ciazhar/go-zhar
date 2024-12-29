package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type StringUserSessionRepository struct {
	redis *redis.Client
}

func NewStringUserSessionRepository(redisClient *redis.Client) *StringUserSessionRepository {
	return &StringUserSessionRepository{
		redis: redisClient,
	}
}

// SetUserSession Set a session token with a TTL
func (r *StringUserSessionRepository) SetUserSession(ctx context.Context, userID string, sessionToken string) error {
	key := fmt.Sprintf("session:%s", userID)
	return r.redis.Set(ctx, key, sessionToken, 24*time.Hour).Err()
}

// GetUserSession Retrieve session token
func (r *StringUserSessionRepository) GetUserSession(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf("session:%s", userID)
	return r.redis.Get(ctx, key).Result()
}
