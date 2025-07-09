package rate_limiter

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func (s *RedisStore) Get(key string) (interface{}, bool) {
	val, err := s.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, false
	}
	return val, true
}

func (s *RedisStore) Set(key string, value interface{}, ttl time.Duration) {
	s.client.Set(context.Background(), key, value, ttl)
}

func (s *RedisStore) Delete(key string) {
	s.client.Del(context.Background(), key)
}
