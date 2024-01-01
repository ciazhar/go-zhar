package rabbitmq_test

import (
	"context"
	"testing"
)

func TestConsumeMessagesChannelPool(t *testing.T) {
	// Consume messages
	channelPool.ConsumeMessages(queueName, HandleMessage)

	// Keep the application running
	forever := make(chan bool)
	<-forever
}

func TestPublishMessageChannelPool(t *testing.T) {
	// Publish the message
	channelPool.PublishMessage(context.Background(), queueName, "Hello, RabbitMQ!")
}
