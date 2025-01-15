package producer

import "time"

const (
	testTopic     = "benchmark_topic"
	messageSize   = 1000 // 1KB
	numMessages   = 100000
	benchmarkTime = 30 * time.Second
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
