package repository_test

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-zhar/examples/redis/crud-testcontainers/internal/repository"
	"github.com/ciazhar/go-zhar/pkg/logger"
	redis2 "github.com/ciazhar/go-zhar/pkg/redis"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestRedisStringRepository(t *testing.T) {
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
	log := logger.Init(logger.Config{
		ConsoleLoggingEnabled: true,
	})

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

	// Test Delete
	err = redisRepo.Delete()
	if err != nil {
		t.Errorf("Failed to delete value from Redis: %v", err)
		return
	}
}

func TestRedisHashRepository(t *testing.T) {
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
	log := logger.Init(logger.Config{
		ConsoleLoggingEnabled: true,
	})

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

	// Test GetHash, SetHash, SetHashTTL, and DeleteHash
	err = redisRepo.SetHash("test_field", "test_value_hash")
	if err != nil {
		t.Errorf("Failed to set hash value in Redis: %v", err)
		return
	}
	value, err := redisRepo.GetHash("test_field")
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

func TestRedisListRepository(t *testing.T) {
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
	log := logger.Init(logger.Config{
		ConsoleLoggingEnabled: true,
	})

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

	// Create - LPUSH to add elements to a list
	list := []string{"value1", "value2", "value3"}
	err = redisRepo.SetList(list)
	if err != nil {
		t.Errorf("Failed to set list value in Redis: %v", err)
		return
	}

	// Read - LRANGE to get elements from the list
	value, err := redisRepo.GetList()
	if err != nil {
		t.Errorf("Failed to get list value from Redis: %v", err)
		return
	}
	expectedList := []string{"value3", "value2", "value1"}
	if !areSlicesEqual(expectedList, value) {
		t.Errorf("Expected list value '%v', got '%s'", list, value)
	}

	// Update - LPUSH to add more elements to the list
	err = redisRepo.SetList([]string{"value0"})
	if err != nil {
		t.Errorf("Failed to set list value in Redis: %v", err)
		return
	}

	// Read - LRANGE to get updated elements from the list
	value, err = redisRepo.GetList()
	if err != nil {
		t.Errorf("Failed to get list value from Redis: %v", err)
		return
	}
	expectedList = []string{"value0", "value3", "value2", "value1"}
	if !areSlicesEqual(value, expectedList) {
		t.Errorf("Expected list value '%v', got '%s'", expectedList, value)
	}

	// Delete - LREM to remove elements from the list
	err = redisRepo.DeleteList("value1")
	if err != nil {
		t.Errorf("Failed to delete list value from Redis: %v", err)
		return
	}

	// Read - LRANGE to get elements after deletion
	value, err = redisRepo.GetList()
	if err != nil {
		t.Errorf("Failed to get list value from Redis: %v", err)
		return
	}
	expectedList = []string{"value0", "value3", "value2"}
	if !areSlicesEqual(value, expectedList) {
		t.Errorf("Expected list value '%v', got '%s'", expectedList, value)
	}
}

func areSlicesEqual(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		fmt.Println("Length of slices are not equal")
		fmt.Println(slice1, slice2)
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}
