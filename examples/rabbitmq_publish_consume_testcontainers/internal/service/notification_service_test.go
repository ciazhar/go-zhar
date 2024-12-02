package service

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/ciazhar/go-start-small/examples/rabbitmq_publish_consume_testcontainers/internal/model"
	rabbitmq2 "github.com/ciazhar/go-start-small/pkg/rabbitmq"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

func TestBasicService(t *testing.T) {
	// Set up RabbitMQ container
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rabbitmqContainer, err := rabbitmq.Run(ctx,
		"rabbitmq:3.12.11-management-alpine",
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

	// Initialize RabbitMQ client
	rabbitMQClient := rabbitmq2.New("test-connection", "guest", "guest", host, port.Port())

	// Create queue
	rabbitMQClient.CreateQueue(model.OrderStatusQueue)
	rabbitMQClient.CreateQueue(model.PaymentReminderQueue)

	// Initialize BasicService
	wg := &sync.WaitGroup{}
	basicService := NewBasicService(ctx, rabbitMQClient, wg)

	// Start RabbitMQ consumer
	basicService.StartRabbitConsumer()

	// Publish messages to RabbitMQ
	basicService.PublishRabbitmq(`{"order_id": "12345","status": "Order Shipped"}`)
	basicService.PublishTTLRabbitmq(`{"order_id": "12345","reminder": "Please complete your payment within 10 minutes."}`)

	<-ctx.Done()
	if err := rabbitmqContainer.Terminate(context.Background()); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	wg.Wait()
}
