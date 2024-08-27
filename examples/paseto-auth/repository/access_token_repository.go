package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// AccessToken represents a JWT access token with associated metadata.
type AccessToken struct {
	UserID      int    `json:"user_id"`
	AuthToken   string `json:"auth_token"`
	GeneratedAt int64  `json:"generated_at"`
	ExpiredAt   int64  `json:"expired_at"`
}

// TokenRepository is an interface for token-related operations.
type TokenRepository interface {
	SaveToken(ctx context.Context, token *AccessToken) error
	GetToken(ctx context.Context, userID int) (*AccessToken, error)
}

// RedisTokenRepository is an implementation of TokenRepository using Redis.
type RedisTokenRepository struct {
	client *redis.Client
}

// NewRedisTokenRepository creates a new instance of RedisTokenRepository.
func NewRedisTokenRepository(client *redis.Client) *RedisTokenRepository {
	return &RedisTokenRepository{
		client: client,
	}
}

// SaveToken saves a token in Redis with a TTL of 5 hours.
func (repo *RedisTokenRepository) SaveToken(ctx context.Context, token *AccessToken) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = repo.client.Set(ctx, tokenKey(token.UserID), data, 5*time.Hour).Err()
	return err
}

// GetToken retrieves a token from Redis by user ID.
func (repo *RedisTokenRepository) GetToken(ctx context.Context, userID int) (*AccessToken, error) {
	data, err := repo.client.Get(ctx, tokenKey(userID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var token AccessToken
	err = json.Unmarshal([]byte(data), &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func tokenKey(userID int) string {
	return "auth_token:" + strconv.Itoa(userID)
}
