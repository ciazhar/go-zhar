package repository

import (
	"github.com/ciazhar/go-zhar/examples/redis/crud-testcontainers/internal/model"
	"github.com/ciazhar/go-zhar/pkg/redis"
	"time"
)

type RedisRepository interface {
	Get() (string, error)
	Set(value string, expiration time.Duration) error
	GetHash(field string) (string, error)
	SetHash(field string, value string) error
	SetHashTTL(field string, value string, ttl time.Duration) error
	DeleteHash(field string) error
}

type redisRepository struct {
	redis redis.Redis
}

func (r redisRepository) Get() (string, error) {
	return r.redis.Get(model.StringKey)
}

func (r redisRepository) Set(value string, expiration time.Duration) error {
	return r.redis.Set(model.StringKey, value, expiration)
}

func (r redisRepository) GetHash(field string) (string, error) {
	return r.redis.GetHash(model.HashKey, field)
}

func (r redisRepository) SetHash(field string, value string) error {
	return r.redis.SetHash(model.HashKey, field, value)
}

func (r redisRepository) SetHashTTL(field string, value string, ttl time.Duration) error {
	return r.redis.SetHashTTL(model.HashKey, field, value, ttl)
}

func (r redisRepository) DeleteHash(field string) error {
	return r.redis.DeleteHash(model.HashKey, field)
}

func NewRedisRepository(redis redis.Redis) RedisRepository {
	return redisRepository{
		redis: redis,
	}
}
