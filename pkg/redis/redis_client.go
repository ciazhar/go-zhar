package redis

import (
	"context"
	"fmt"

	"github.com/ciazhar/go-zhar/pkg/bootstrap"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/redis/go-redis/v9"
)

func InitRedis(ctx context.Context, host string, port int, password string) (*bootstrap.ClientService, *redis.Client) {
	var (
		log = logger.FromContext(ctx).With().Str("host", host).Int("port", port).Logger()
	)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       0, // use default DB
	})
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal().Err(err).Msg("failed to ping redis")
	}
	return bootstrap.NewClientService(ctx, "redis", func() error {
		client.Close()
		return nil
	}), client
}
