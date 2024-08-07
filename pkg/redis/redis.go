package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type Redis struct {
	rdb    *redis.Client
	logger *logger.Logger
}

func Init(host string, port int, password string, logger *logger.Logger) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       0, // use default DB
	})

	status := rdb.Ping(context.Background())
	if status.Err() != nil {
		logger.Fatalf("Error connecting to redis: %s", status.Err())
	}

	logger.Info("Redis client initialized successfully.")

	return &Redis{
		rdb:    rdb,
		logger: logger,
	}
}

func (r *Redis) Get(key string) (string, error) {
	val, err := r.rdb.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", fmt.Errorf("%s: %s", "Error getting value from redis", err)
	}
	return val, nil
}

// Set sets a key-value pair in Redis with an expiration time.
//
// Parameters:
//   - key: the key to set in Redis.
//   - value: the value corresponding to the key.
//   - expiration: the duration for which the key-value pair should be stored. Set 0 if the value should not expire.
//
// Returns an error if there was an issue setting the value.
func (r *Redis) Set(key string, value string, expiration time.Duration) error {
	_, err := r.rdb.Set(context.Background(), key, value, expiration).Result()
	if err != nil {
		return fmt.Errorf("%s: %s", "Error setting value in redis", err)
	}
	return nil
}

func (r *Redis) Delete(key string) error {
	_, err := r.rdb.Del(context.Background(), key).Result()
	if err != nil {
		return fmt.Errorf("%s: %s", "Error deleting value from redis", err)
	}
	return nil
}

func (r *Redis) GetHash(key string, field string) (string, error) {
	val, err := r.rdb.HGet(context.Background(), key, field).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", fmt.Errorf("%s: %s", "Error getting value from redis", err)
	}
	return val, nil
}

func (r *Redis) SetHash(key string, field string, value string) error {
	_, err := r.rdb.HSet(context.Background(), key, field, value).Result()
	if err != nil {
		return fmt.Errorf("%s: %s", "Error setting value in redis", err)
	}
	return nil
}

func (r *Redis) SetHashTTL(key string, field string, value string, ttl time.Duration) error {
	_, err := r.rdb.HSet(context.Background(), key, field, value).Result()
	if err != nil {
		return fmt.Errorf("%s: %s", "Error setting value in redis", err)
	}

	err = r.rdb.Expire(context.Background(), key, ttl).Err()
	if err != nil {
		return fmt.Errorf("%s: %s", "Error setting TTL on hash", err)
	}

	return nil
}

func (r *Redis) DeleteHash(key string, field string) error {
	_, err := r.rdb.HDel(context.Background(), key, field).Result()
	if err != nil {
		return fmt.Errorf("%s: %s", "Error deleting value from redis", err)
	}
	return nil
}

func (r *Redis) GetList(key string) ([]string, error) {
	val, err := r.rdb.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %s", "Error getting value from redis", err)
	}
	return val, nil
}

func (r *Redis) SetList(key string, list []string) error {
	err := r.rdb.LPush(context.Background(), key, list).Err()
	if err != nil {
		return fmt.Errorf("%s: %s", "Error setting value in redis", err)
	}
	return nil
}

func (r *Redis) DeleteList(key string, value string) error {
	_, err := r.rdb.LRem(context.Background(), key, 1, value).Result()
	if err != nil {
		return fmt.Errorf("%s: %s", "Error deleting value from redis", err)
	}
	return nil
}

func (r *Redis) Close() {
	defer func() {
		err := r.rdb.Close()
		if err != nil {
			log.Fatalf("%s: %s", "Error closing redis", err)
		}
	}()

}
