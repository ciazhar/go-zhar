package pkg

import "github.com/go-redis/redis/v8"

var RedisClient *redis.Client

func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6377",
	})
}
