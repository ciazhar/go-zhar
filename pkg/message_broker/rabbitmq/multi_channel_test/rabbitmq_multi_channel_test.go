package multi_channel_test

import (
	"fmt"
	"github.com/ciazhar/zhar/pkg/message_broker/rabbitmq/basic"
	rabbitmq2 "github.com/ciazhar/zhar/pkg/message_broker/rabbitmq/multi_channel_test"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"testing"
)

var (
	pool      *rabbitmq2.ChannelPool
	ch        *amqp.Channel
	queueName = "test_queue"
	err       error
)

func TestMain(m *testing.M) {
	// Set up the shared connection and channel before running tests
	pool = rabbitmq2.NewChannelPool("guest", "guest", "localhost", "5672", 5)
	ch, err = pool.Get()
	if err != nil {
		log.Fatalf("Error getting connection from pool: %s", err)
	}
	defer pool.Put(ch)

	// Create the queue
	basic.CreateQueue(ch, queueName)

	// Run the tests
	exitCode := m.Run()

	// Exit with the status code from tests
	os.Exit(exitCode)
}

func TestConsumeMessages(t *testing.T) {

	// Consume messages
	basic.ConsumeMessages(ch, queueName, HandleMessage)

	// Keep the application running
	forever := make(chan bool)
	<-forever
}

func HandleMessage(string2 string) {
	fmt.Printf("Received a message: %s\n", string2)
}

func TestPublishMessage(t *testing.T) {
	basic.PublishMessage(ch, queueName, "Hello, RabbitMQ!")
}
