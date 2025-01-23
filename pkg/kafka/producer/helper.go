package producer

import "errors"

const (
	testTopic   = "benchmark_topic"
	messageSize = 1000 // 1KB
	numMessages = 100000
)

var brokers = []string{"localhost:9092", "localhost:9093", "localhost:9094"}

// generateMessage creates a message of specified size
func generateMessage(size int) string {
	message := make([]byte, size)
	for i := 0; i < size; i++ {
		message[i] = 'a' + byte(i%26)
	}
	return string(message)
}

var ErrTxnInProgress = errors.New("transaction already in progress")
var ErrProducerClosed = errors.New("producer is closed")
var ErrInvalidMessage = errors.New("invalid message")
var ErrEmptyBatch = errors.New("empty message batch")
