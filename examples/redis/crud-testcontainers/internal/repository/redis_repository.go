package repository

import (
	"github.com/ciazhar/go-zhar/examples/redis/crud-testcontainers/internal/model"
	"github.com/ciazhar/go-zhar/pkg/redis"
	"time"
)

type RedisRepository struct {
	redis *redis.Redis
}

func (r *RedisRepository) Get() (string, error) {
	return r.redis.Get(model.StringKey)
}

func (r *RedisRepository) Set(value string, expiration time.Duration) error {
	return r.redis.Set(model.StringKey, value, expiration)
}

func (r *RedisRepository) Delete() error {
	return r.redis.Delete(model.StringKey)
}

func (r *RedisRepository) GetHash(field string) (string, error) {
	return r.redis.GetHash(model.HashKey, field)
}

func (r *RedisRepository) SetHash(field string, value string) error {
	return r.redis.SetHash(model.HashKey, field, value)
}

func (r *RedisRepository) SetHashTTL(field string, value string, ttl time.Duration) error {
	return r.redis.SetHashTTL(model.HashKey, field, value, ttl)
}

func (r *RedisRepository) DeleteHash(field string) error {
	return r.redis.DeleteHash(model.HashKey, field)
}

func (r *RedisRepository) GetList() ([]string, error) {
	return r.redis.GetList(model.ListKey)
}

func (r *RedisRepository) SetList(list []string) error {
	return r.redis.SetList(model.ListKey, list)
}

func (r *RedisRepository) DeleteList(value string) error {
	return r.redis.DeleteList(model.ListKey, value)
}

func NewRedisRepository(redis *redis.Redis) RedisRepository {
	return RedisRepository{
		redis: redis,
	}
}
