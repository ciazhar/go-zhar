package producer

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkCustomPartitionerProducer_Throughput(b *testing.B) {

	producer, err := NewCustomPartitionerProducer(brokers)
	if err != nil {
		b.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	message := generateMessage(messageSize)

	start := time.Now()
	for i := 0; i < numMessages; i++ {
		key := fmt.Sprintf("key-%d", i)
		err := producer.SendMessage(testTopic, key, message)
		if err != nil {
			b.Fatalf("Failed to send message: %v", err)
		}
	}

	duration := time.Since(start)

	messagesPerSecond := float64(numMessages) / duration.Seconds()
	b.ReportMetric(messagesPerSecond, "msgs/sec")
	b.ReportMetric(float64(messageSize*numMessages)/duration.Seconds(), "bytes/sec")
}

// Add more specific benchmarks for different partition scenarios
func BenchmarkCustomPartitionerProducer_PartitionDistribution(b *testing.B) {

	producer, err := NewCustomPartitionerProducer(brokers)
	if err != nil {
		b.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	message := generateMessage(100) // Smaller message size for partition testing

	// Test with different key patterns
	patterns := []struct {
		name    string
		keyFunc func(i int) string
	}{
		{"Sequential", func(i int) string { return fmt.Sprintf("key-%d", i) }},
		{"Constant", func(i int) string { return "constant-key" }},
		{"Random", func(i int) string { return fmt.Sprintf("key-%d", rand.Intn(1000)) }},
	}

	for _, pattern := range patterns {
		b.Run(pattern.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				key := pattern.keyFunc(i)
				err := producer.SendMessage(testTopic, key, message)
				if err != nil {
					b.Fatalf("Failed to send message: %v", err)
				}
			}
		})
	}
}
