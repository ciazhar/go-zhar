package rabbitmq_test

import (
	"context"
	"testing"
)

func TestConsumeMessagesConnectionPool(t *testing.T) {

	// Consume messages
	connectionPool.ConsumeMessages(queueName, HandleMessage)

	// Keep the application running
	forever := make(chan bool)
	<-forever
}

func TestPublishMessageConnectionPool(t *testing.T) {
	connectionPool.PublishMessage(context.Background(), queueName, "Hello, RabbitMQ!")
}
