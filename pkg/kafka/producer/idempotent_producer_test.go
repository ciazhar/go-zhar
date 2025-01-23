package producer

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkIdempotentProducer_Throughput(b *testing.B) {

	producer, err := NewIdempotentProducer(brokers)
	if err != nil {
		b.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	message := generateMessage(messageSize)

	start := time.Now()
	for i := 0; i < numMessages; i++ {
		key := fmt.Sprintf("key-%d", i)
		_, _, err := producer.SendMessage(testTopic, key, message)
		if err != nil {
			b.Fatalf("Failed to send message: %v", err)
		}
	}

	duration := time.Since(start)

	messagesPerSecond := float64(numMessages) / duration.Seconds()
	b.ReportMetric(messagesPerSecond, "msgs/sec")
	b.ReportMetric(float64(messageSize*numMessages)/duration.Seconds(), "bytes/sec")
}

// Benchmark with different message sizes
func BenchmarkIdempotentProducer_MessageSizes(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size-%dB", size), func(b *testing.B) {
			producer, err := NewIdempotentProducer(brokers)
			if err != nil {
				b.Fatalf("Failed to create producer: %v", err)
			}
			defer producer.Close()

			message := generateMessage(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := fmt.Sprintf("key-%d", i)
				_, _, err := producer.SendMessage(testTopic, key, message)
				if err != nil {
					b.Fatalf("Failed to send message: %v", err)
				}
			}
		})
	}
}

// Benchmark with parallel sends
func BenchmarkIdempotentProducer_Parallel(b *testing.B) {
	producer, err := NewIdempotentProducer(brokers)
	if err != nil {
		b.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	message := generateMessage(messageSize)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			counter++
			key := fmt.Sprintf("key-%d", counter)
			_, _, err := producer.SendMessage(testTopic, key, message)
			if err != nil {
				b.Fatalf("Failed to send message: %v", err)
			}
		}
	})
}
