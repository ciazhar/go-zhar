package redis

import (
	"context"
	"fmt"

	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/go-redis/redis/v8"
)

func InitRedis(host string, port int, password string) *redis.Client {

	logger.LogInfo(context.Background(), "Connecting to Redis", map[string]interface{}{
		"host": host,
		"port": port,
	})

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       0, // use default DB
	})

	status := rdb.Ping(context.Background())
	if status.Err() != nil {
		logger.LogFatal(context.Background(), status.Err(), "Failed to connect to Redis", map[string]interface{}{
			"status": status.Val(),
		})
	}

	logger.LogInfo(context.Background(), "Connected to Redis", nil)

	return rdb
}
