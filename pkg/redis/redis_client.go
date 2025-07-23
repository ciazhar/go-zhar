package redis

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-start-small/pkg/bootstrap"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/redis/go-redis/v9"
)

func InitRedis(ctx context.Context, host string, port int, password string) (*bootstrap.ClientService, *redis.Client, error) {

	var (
		log = logger.FromContext(ctx).With().Str("host", host).Int("port", port).Logger()
	)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       0, // use default DB
	})
	if err := client.Ping(ctx).Err(); err != nil {
		log.Err(err).Send()
		return nil, nil, err
	}
	return bootstrap.NewClientService("redis", func() error {
		client.Close()
		return nil
	}), client, nil
}
