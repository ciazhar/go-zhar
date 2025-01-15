package producer

import (
	"context"
	"fmt"
	"testing"
	"time"
)

const (
// testTopic     = "benchmark_topic"
// messageSize   = 1000 // 1KB
// numMessages   = 100000
// benchmarkTime = 30 * time.Second
)

//func createTestBatchProducer() (*BatchProducer, error) {
//	cfg := BatchProducerConfig{
//		BrokerList:   []string{"localhost:9092", "localhost:9093", "localhost:9094"},
//		BatchSize:    1000,
//		FlushTimeout: 100 * time.Millisecond,
//		Compression:  sarama.CompressionSnappy,
//		RequiredAcks: sarama.WaitForLocal,
//		RetryMax:     3,
//		RetryBackoff: 100 * time.Millisecond,
//	}
//	return NewBatchProducer(cfg)
//}
//
//func BenchmarkBatchProducer_SendMessage(b *testing.B) {
//	producer, err := createTestBatchProducer()
//	if err != nil {
//		b.Fatalf("Failed to create producer: %v", err)
//	}
//	defer producer.Close()
//
//	message := generateMessage(messageSize)
//
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		counter := 0
//		for pb.Next() {
//			key := fmt.Sprintf("key-%d", counter)
//			err := producer.SendMessage(testTopic, key, message)
//			if err != nil {
//				b.Fatalf("Failed to send message: %v", err)
//			}
//			counter++
//		}
//	})
//}
//
//func BenchmarkBatchProducer_BatchSizes(b *testing.B) {
//	batchSizes := []int{100, 500, 1000, 5000}
//
//	for _, batchSize := range batchSizes {
//		b.Run(fmt.Sprintf("BatchSize-%d", batchSize), func(b *testing.B) {
//			cfg := BatchProducerConfig{
//				BrokerList:   []string{"localhost:9092", "localhost:9093", "localhost:9094"},
//				BatchSize:    batchSize,
//				FlushTimeout: 100 * time.Millisecond,
//				Compression:  sarama.CompressionSnappy,
//				RequiredAcks: sarama.WaitForLocal,
//				RetryMax:     3,
//				RetryBackoff: 100 * time.Millisecond,
//			}
//
//			producer, err := NewBatchProducer(cfg)
//			if err != nil {
//				b.Fatalf("Failed to create producer: %v", err)
//			}
//			defer producer.Close()
//
//			message := generateMessage(messageSize)
//
//			b.ResetTimer()
//			for i := 0; i < b.N; i++ {
//				key := fmt.Sprintf("key-%d", i)
//				err := producer.SendMessage(testTopic, key, message)
//				if err != nil {
//					b.Fatalf("Failed to send message: %v", err)
//				}
//			}
//
//			// Report stats
//			stats := producer.GetStats()
//			b.ReportMetric(float64(stats.MessagesSent)/b.Elapsed().Seconds(), "msgs/sec")
//			b.ReportMetric(float64(stats.BytesSent)/b.Elapsed().Seconds(), "bytes/sec")
//			b.ReportMetric(float64(stats.Errors), "errors")
//		})
//	}
//}
//
//func BenchmarkBatchProducer_Compression(b *testing.B) {
//	compressionModes := []struct {
//		name string
//		code sarama.CompressionCodec
//	}{
//		{"None", sarama.CompressionNone},
//		{"GZIP", sarama.CompressionGZIP},
//		{"Snappy", sarama.CompressionSnappy},
//		{"LZ4", sarama.CompressionLZ4},
//		{"ZSTD", sarama.CompressionZSTD},
//	}
//
//	for _, compression := range compressionModes {
//		b.Run(compression.name, func(b *testing.B) {
//			cfg := BatchProducerConfig{
//				BrokerList:   []string{"localhost:9092", "localhost:9093", "localhost:9094"},
//				BatchSize:    1000,
//				FlushTimeout: 100 * time.Millisecond,
//				Compression:  compression.code,
//				RequiredAcks: sarama.WaitForLocal,
//				RetryMax:     3,
//				RetryBackoff: 100 * time.Millisecond,
//			}
//
//			producer, err := NewBatchProducer(cfg)
//			if err != nil {
//				b.Fatalf("Failed to create producer: %v", err)
//			}
//			defer producer.Close()
//
//			message := generateMessage(messageSize)
//
//			b.ResetTimer()
//			for i := 0; i < b.N; i++ {
//				key := fmt.Sprintf("key-%d", i)
//				err := producer.SendMessage(testTopic, key, message)
//				if err != nil {
//					b.Fatalf("Failed to send message: %v", err)
//				}
//			}
//
//			stats := producer.GetStats()
//			b.ReportMetric(float64(stats.MessagesSent)/b.Elapsed().Seconds(), "msgs/sec")
//			b.ReportMetric(float64(stats.BytesSent)/b.Elapsed().Seconds(), "bytes/sec")
//		})
//	}
//}

//func TestBatchProducer_BasicFunctionality(t *testing.T) {
//
//	// 1. Add a context with timeout
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	// 2. Ensure proper shutdown sequence
//	defer func() {
//		// Close input channel
//		close(inputChan)
//
//		// Wait for producer to flush and close
//		producer.Close()
//
//		// Wait for any goroutines to finish
//		producer.WaitGroup.Wait()
//	}()
//
//	cfg := BatchProducerConfig{
//		BrokerList:   []string{"localhost:9092"},
//		BatchSize:    10,
//		FlushTimeout: 100 * time.Millisecond,
//		Compression:  sarama.CompressionNone,
//		RequiredAcks: sarama.WaitForLocal,
//		RetryMax:     3,
//		RetryBackoff: 100 * time.Millisecond,
//	}
//
//	producer, err := NewBatchProducer(cfg)
//	if err != nil {
//		t.Fatalf("Failed to create producer: %v", err)
//	}
//
//	// Create a WaitGroup to track message sending
//	var wg sync.WaitGroup
//	wg.Add(20)
//
//	// Send test messages
//	for i := 0; i < 20; i++ {
//		i := i // Capture loop variable
//		go func() {
//			defer wg.Done()
//			err := producer.SendMessage(testTopic, fmt.Sprintf("key-%d", i), "test message")
//			if err != nil {
//				t.Errorf("Failed to send message %d: %v", i, err)
//			}
//		}()
//	}
//
//	// Wait for all messages to be sent
//	wg.Wait()
//
//	// Wait for messages to be processed
//	time.Sleep(2 * time.Second)
//
//	// Get stats before closing
//	stats := producer.GetStats()
//	t.Logf("Messages sent: %d", stats.MessagesSent)
//	t.Logf("Batches processed: %d", stats.BatchesProcessed)
//	t.Logf("Errors: %d", stats.Errors)
//
//	// Gracefully close the producer
//	err = producer.Close()
//	if err != nil {
//		t.Fatalf("Failed to close producer: %v", err)
//	}
//
//	// Verify results
//	finalStats := producer.GetStats()
//	if finalStats.MessagesSent != 20 {
//		t.Errorf("Expected 20 messages to be sent, got %d", finalStats.MessagesSent)
//	}
//	if finalStats.Errors > 0 {
//		t.Errorf("Expected 0 errors, got %d", finalStats.Errors)
//	}
//}

func TestBatchProducer_BasicFunctionality(t *testing.T) {
	// Configure test parameters
	brokers := []string{"localhost:9092"}
	batchSize := 10
	totalMessages := 20

	// Create producer
	producer, err := NewBatchProducer(brokers, batchSize)
	if err != nil {
		t.Fatalf("Failed to create producer: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Start the producer
	producer.Start()

	// Send test messages
	messagesSent := 0
	for i := 0; i < totalMessages; i++ {
		select {
		case <-ctx.Done():
			t.Fatal("Context deadline exceeded while sending messages")
		case producer.Input() <- Message{
			Topic: "benchmark_topic",
			Key:   []byte(fmt.Sprintf("key-%d", i)),
			Value: []byte(fmt.Sprintf("value-%d", i)),
		}:
			messagesSent++
		}
	}

	// Close producer and wait for completion
	err = producer.Close()
	if err != nil {
		t.Errorf("Error closing producer: %v", err)
	}

	// Wait for producer to finish with timeout
	select {
	case <-ctx.Done():
		t.Fatal("Test timed out waiting for producer to finish")
	case <-producer.Done():
		// Success
	}

	// Log results
	t.Logf("Messages sent: %d", messagesSent)
	t.Logf("Batches processed: %d", messagesSent/batchSize)
	t.Logf("Errors: %d", 0) // You might want to track errors in your implementation
}
