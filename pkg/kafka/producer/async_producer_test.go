package producer

import (
	"fmt"
	"testing"
	"time"
)

// BenchmarkAsyncProducer_Throughput measures messages/second
func BenchmarkAsyncProducer_Throughput(b *testing.B) {
	producer, err := NewAsyncProducer(brokers)
	if err != nil {
		b.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	message := generateMessage(messageSize)

	// Measure time taken to send numMessages
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

// BenchmarkAsyncProducer_DifferentSizes benchmarks different message sizes
func BenchmarkAsyncProducer_DifferentSizes(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000} // 100B, 1KB, 10KB, 100KB

	for _, size := range sizes {
		b.Run(fmt.Sprintf("MessageSize-%dB", size), func(b *testing.B) {
			brokers := []string{"localhost:9092", "localhost:9093", "localhost:9094"}
			producer, err := NewAsyncProducer(brokers)
			if err != nil {
				b.Fatalf("Failed to create producer: %v", err)
			}
			defer producer.Close()

			message := generateMessage(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := fmt.Sprintf("key-%d", i)
				err := producer.SendMessage(testTopic, key, message)
				if err != nil {
					b.Fatalf("Failed to send message: %v", err)
				}
			}
		})
	}
}
