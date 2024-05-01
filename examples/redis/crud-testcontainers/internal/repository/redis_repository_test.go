package repository_test

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/redis/crud-testcontainers/internal/repository"
	"github.com/ciazhar/go-zhar/pkg/logger"
	redis2 "github.com/ciazhar/go-zhar/pkg/redis"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestRedisRepository(t *testing.T) {
	// Set up Redis container
	redisPassword := "your_redis_password"
	ctx := context.Background()
	redisContainer, err := redis.RunContainer(ctx,
		testcontainers.WithImage("docker.io/redis:7"),
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	// Clean up the container
	defer func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// Get Redis container host and port
	host, err := redisContainer.Host(ctx)
	if err != nil {
		t.Errorf("Failed to get Redis container host: %v", err)
		return
	}
	port, err := redisContainer.MappedPort(ctx, "6379/tcp")
	if err != nil {
		t.Errorf("Failed to get Redis container port: %v", err)
		return
	}

	// Logger configuration
	log := logger.Init()

	// Initialize Redis client
	r := redis2.Init(
		host,
		port.Int(),
		redisPassword,
		log,
	)
	defer r.Close()

	// Initialize Redis repository
	redisRepo := repository.NewRedisRepository(r)

	// Test Get and Set
	err = redisRepo.Set("test_value", 0)
	if err != nil {
		t.Errorf("Failed to set value in Redis: %v", err)
		return
	}
	value, err := redisRepo.Get()
	if err != nil {
		t.Errorf("Failed to get value from Redis: %v", err)
		return
	}
	if value != "test_value" {
		t.Errorf("Expected value 'test_value', got '%s'", value)
	}

	// Test GetHash, SetHash, SetHashTTL, and DeleteHash
	err = redisRepo.SetHash("test_field", "test_value_hash")
	if err != nil {
		t.Errorf("Failed to set hash value in Redis: %v", err)
		return
	}
	value, err = redisRepo.GetHash("test_field")
	if err != nil {
		t.Errorf("Failed to get hash value from Redis: %v", err)
		return
	}
	if value != "test_value_hash" {
		t.Errorf("Expected hash value 'test_value_hash', got '%s'", value)
	}

	err = redisRepo.SetHashTTL("test_field_ttl", "test_value_ttl", 5*time.Second)
	if err != nil {
		t.Errorf("Failed to set hash value with TTL in Redis: %v", err)
		return
	}

	time.Sleep(6 * time.Second)

	value, err = redisRepo.GetHash("test_field_ttl")
	if value != "" {
		t.Error("Expected empty hash value, got value")
	}

	err = redisRepo.DeleteHash("test_field")
	if err != nil {
		t.Errorf("Failed to delete hash value from Redis: %v", err)
		return
	}

	value, err = redisRepo.GetHash("test_field")
	if value != "" {
		t.Error("Expected empty hash value, got value")
	}
}
