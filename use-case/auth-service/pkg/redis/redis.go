package redis

import (
	"github.com/go-redis/redis/v8"
	"os"
)

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
	})
}
