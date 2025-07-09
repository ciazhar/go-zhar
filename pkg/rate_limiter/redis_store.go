package rate_limiter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func (r *RedisStore) Get(key string, out interface{}) (bool, error) {
	result, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil // key not found
		}
		return false, fmt.Errorf("redis get error: %w", err)
	}

	if err := json.Unmarshal([]byte(result), out); err != nil {
		return false, fmt.Errorf("unmarshal error: %w", err)
	}

	return true, nil
}

func (r *RedisStore) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return r.client.Set(context.Background(), key, data, ttl).Err()
}

func (s *RedisStore) Delete(key string) {
	s.client.Del(context.Background(), key)
}

func (s *RedisStore) Type() StorageType {
	return Redis
}
