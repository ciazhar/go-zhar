package repository

import (
	"github.com/ciazhar/go-zhar/examples/cache/redis/basic/internal/basic/model"
	"github.com/ciazhar/go-zhar/pkg/cache/redis"
)

type RedisRepository interface {
	GetBasicHash(key string) (string, error)
	SetBasicHash(key string, value string) error
}

type redisRepository struct {
	redis redis.Redis
}

func (r redisRepository) GetBasicHash(key string) (string, error) {
	return r.redis.GetHash(model.BasicKey, key)
}

func (r redisRepository) SetBasicHash(key string, value string) error {
	return r.redis.SetHash(model.BasicKey, key, value)
}

func NewRedisRepositoryParams(redis redis.Redis) RedisRepository {
	return redisRepository{
		redis: redis,
	}
}
