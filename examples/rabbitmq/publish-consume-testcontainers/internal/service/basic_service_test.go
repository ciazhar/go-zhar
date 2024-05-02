package service

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/rabbitmq/publish-consume-testcontainers/internal/model"
	logger2 "github.com/ciazhar/go-zhar/pkg/logger"
	rabbitmq2 "github.com/ciazhar/go-zhar/pkg/rabbitmq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"log"
	"sync"
	"testing"
	"time"
)

func TestBasicService(t *testing.T) {
	// Set up RabbitMQ container
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rabbitmqContainer, err := rabbitmq.RunContainer(ctx,
		testcontainers.WithImage("rabbitmq:3.12.11-management-alpine"),
		rabbitmq.WithAdminUsername("admin"),
		rabbitmq.WithAdminPassword("password"),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// Get RabbitMQ container host and port
	host, err := rabbitmqContainer.Host(ctx)
	if err != nil {
		t.Errorf("Failed to get RabbitMQ container host: %v", err)
		return
	}
	port, err := rabbitmqContainer.MappedPort(ctx, "5672/tcp")
	if err != nil {
		t.Errorf("Failed to get RabbitMQ container port: %v", err)
		return
	}

	// Set up logger
	log := logger2.Init()

	// Initialize RabbitMQ client
	rabbitMQClient := rabbitmq2.New("test-connection", "guest", "guest", host, port.Port(), log)

	// Create queue
	rabbitMQClient.CreateQueue(model.QueueBasic)

	// Initialize BasicService
	wg := &sync.WaitGroup{}
	basicService := NewBasicService(ctx, rabbitMQClient, wg, log)

	// Start RabbitMQ consumer
	basicService.StartRabbitConsumer()

	// Publish messages to RabbitMQ
	basicService.PublishRabbitmq("test_message_1")
	basicService.PublishTTLRabbitmq("test_message_2")

	select {
	case <-ctx.Done():
		if err := rabbitmqContainer.Terminate(context.Background()); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}

	wg.Wait()
}
