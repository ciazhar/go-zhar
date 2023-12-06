package multi_connection_test

import (
	"fmt"
	error2 "github.com/ciazhar/zhar/pkg/error"
	"github.com/ciazhar/zhar/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"testing"
)

var (
	pool      *rabbitmq.ConnectionPool
	ch        *amqp.Channel
	queueName = "test_queue"
)

func TestMain(m *testing.M) {
	// Set up the shared connection and channel before running tests
	pool = rabbitmq.NewConnectionPool("guest", "guest", "localhost", "5672", 5)
	conn, err := pool.Get()
	if err != nil {
		log.Fatalf("Error getting connection from pool: %s", err)
	}
	defer pool.Put(conn)

	ch, err = conn.Channel()
	error2.FailOnError(err, "Failed to open a channel")
	fmt.Println("Channel created")

	// Create the queue
	rabbitmq.CreateQueue(ch, queueName)

	// Run the tests
	exitCode := m.Run()

	// Exit with the status code from tests
	os.Exit(exitCode)
}

func TestConsumeMessages(t *testing.T) {

	// Consume messages
	rabbitmq.ConsumeMessages(ch, queueName, HandleMessage)

	// Keep the application running
	forever := make(chan bool)
	<-forever
}

func HandleMessage(string2 string) {
	fmt.Printf("Received a message: %s\n", string2)
}

func TestPublishMessage(t *testing.T) {
	rabbitmq.PublishMessage(ch, queueName, "Hello, RabbitMQ!")
}
