package rabbitmq_test

import (
	"github.com/ciazhar/zhar/pkg/message_broker/rabbitmq"
	"os"
	"testing"
)

var (
	queueName      = "test_queue"
	rabbit         *rabbitmq.RabbitMQ
	channelPool    *rabbitmq.ChannelPool
	connectionPool *rabbitmq.ConnectionPool
)

func TestMain(m *testing.M) {

	// Set up the shared connection and channel before running tests
	rabbit = rabbitmq.New("guest", "guest", "localhost", "5672")
	defer rabbit.Close()

	// Set up the shared connection and channel before running tests
	channelPool = rabbitmq.NewChannelPool("guest", "guest", "localhost", "5672", 5)
	defer channelPool.Close()

	// Set up the shared connection and channel before running tests
	connectionPool = rabbitmq.NewConnectionPool("guest", "guest", "localhost", "5672", 5)
	defer connectionPool.Close()

	// Create the queue
	rabbit.CreateQueue(queueName)

	// Run the tests
	exitCode := m.Run()

	// Exit with the status code from tests
	os.Exit(exitCode)
}
