package rabbitmq_test

import (
	"context"
	"fmt"
	"testing"
)

func TestConsumeMessages(t *testing.T) {
	// Consume messages
	rabbit.ConsumeMessages(queueName, HandleMessage)

	// Keep the application running
	forever := make(chan bool)
	<-forever
}

func HandleMessage(string2 string) {
	fmt.Printf("Received a message: %s\n", string2)
}

func TestPublishMessage(t *testing.T) {
	// Create the queue
	rabbit.CreateQueue(queueName)

	// Publish the message
	rabbit.PublishMessage(context.Background(), queueName, "Hello, RabbitMQ!")
}
