package producer

import (
	"fmt"
	"testing"
	"time"
)

// BenchmarkAsyncProducer_SendMessage benchmarks the sending of messages
func BenchmarkAsyncProducer_SendMessage(b *testing.B) {
	// Connect to your Kafka brokers
	producer, err := NewAsyncProducer(brokers)
	if err != nil {
		b.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	message := generateMessage(messageSize)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", counter)
			err = producer.SendMessage(testTopic, key, message)
			if err != nil {
				b.Fatalf("Failed to send message: %v", err)
			}
			counter++
		}
	})
}

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

// BenchmarkAsyncProducer_TimeBound runs for a fixed duration
func BenchmarkAsyncProducer_TimeBound(b *testing.B) {
	producer, err := NewAsyncProducer(brokers)
	if err != nil {
		b.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	message := generateMessage(messageSize)

	done := make(chan struct{})
	messageCount := 0

	// Start timer
	start := time.Now()

	// Send messages
	go func() {
		for time.Since(start) < benchmarkTime {
			key := fmt.Sprintf("key-%d", messageCount)
			err := producer.SendMessage(testTopic, key, message)
			if err != nil {
				b.Errorf("Failed to send message: %v", err)
				return
			}
			messageCount++
		}
		close(done)
	}()

	<-done
	duration := time.Since(start)

	messagesPerSecond := float64(messageCount) / duration.Seconds()
	b.ReportMetric(messagesPerSecond, "msgs/sec")
	b.ReportMetric(float64(messageSize*messageCount)/duration.Seconds(), "bytes/sec")
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
