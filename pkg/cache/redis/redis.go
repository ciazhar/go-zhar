package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

type Redis struct {
	rdb *redis.Client
}

// Init initializes a Redis client.
//
// Parameters:
// - host: a string representing the host address.
// - port: an integer representing the port number.
// - password: a string representing the password for authentication.
//
// Returns a Redis struct.
func Init(host string, port int, password string) Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       0, // use default DB
	})

	log.Println("Redis client initialized successfully.")

	return Redis{
		rdb: rdb,
	}
}

// GetHash retrieves the value of a field in a Redis hash.
//
// Parameters:
// - key: the key of the hash.
// - field: the field within the hash.
//
// Returns:
// - string: the value of the field.
// - error: an error, if any.
func (r Redis) GetHash(key string, field string) (string, error) {
	val, err := r.rdb.HGet(context.Background(), key, field).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// SetHash sets a hash field to a value in Redis.
//
// Parameters:
// - key: the key of the hash.
// - field: the field within the hash.
// - value: the value to set.
// Return type: error.
func (r Redis) SetHash(key string, field string, value string) error {
	_, err := r.rdb.HSet(context.Background(), key, field, value).Result()
	if err != nil {
		return err
	}
	return nil
}

// Close closes the Redis connection.
//
// This function does not take any parameters.
// It does not return any values.
func (r Redis) Close() {
	defer func() {
		err := r.rdb.Close()
		if err != nil {
			log.Fatalf("%s: %s", "Error closing redis", err)
		}
	}()

}
