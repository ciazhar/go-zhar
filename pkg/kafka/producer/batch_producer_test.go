package producer

import (
	"fmt"
	"github.com/IBM/sarama"
	"testing"
	"time"
)

// BenchmarkBatchProducer_SendMessage benchmarks the sending of messages
func BenchmarkBatchProducer_SendMessage(b *testing.B) {
	config := ProducerConfig{
		BatchSize:   100,
		Compression: sarama.CompressionNone,
	}

	producer, err := NewBatchProducer(brokers, config)
	if err != nil {
		b.Fatalf("Failed to create producer: %v", err)
	}
	producer.Start()
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

func BenchmarkBatchProducer_BatchSizes(b *testing.B) {
	batchSizes := []int{1, 10, 100, 1000}

	message := generateMessage(messageSize)

	for _, batchSize := range batchSizes {
		b.Run(fmt.Sprintf("BatchSize-%d", batchSize), func(b *testing.B) {
			config := ProducerConfig{
				BatchSize:   batchSize,
				Compression: sarama.CompressionNone,
			}

			producer, err := NewBatchProducer(brokers, config)
			if err != nil {
				b.Fatalf("Failed to create producer: %v", err)
			}
			producer.Start()
			defer producer.Close()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := fmt.Sprintf("key-%d", i)
				err = producer.SendMessage(testTopic, key, message)
				if err != nil {
					b.Fatalf("Failed to send message: %v", err)
				}
			}
		})
	}
}

func BenchmarkBatchProducer_Compression(b *testing.B) {
	compressionCodes := []sarama.CompressionCodec{
		sarama.CompressionNone,
		sarama.CompressionGZIP,
		sarama.CompressionSnappy,
		sarama.CompressionLZ4,
		sarama.CompressionZSTD,
	}

	for _, compression := range compressionCodes {
		b.Run(fmt.Sprintf("Compression-%d", compression), func(b *testing.B) {
			config := ProducerConfig{
				BatchSize:   100,
				Compression: compression,
			}

			producer, err := NewBatchProducer(brokers, config)
			if err != nil {
				b.Fatalf("Failed to create producer: %v", err)
			}
			producer.Start()
			defer producer.Close()

			message := generateMessage(messageSize)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := fmt.Sprintf("key-%d", i)
				err = producer.SendMessage(testTopic, key, message)
				if err != nil {
					b.Fatalf("Failed to send message: %v", err)
				}
			}
		})
	}
}

func BenchmarkBatchProducer_Throughput(b *testing.B) {
	configs := []struct {
		name      string
		batchSize int
	}{
		{"BatchSize-100", 100},
		{"BatchSize-500", 500},
		{"BatchSize-1000", 1000},
	}

	message := generateMessage(messageSize)

	for _, cfg := range configs {
		b.Run(cfg.name, func(b *testing.B) {
			config := ProducerConfig{
				BatchSize:   cfg.batchSize,
				Compression: sarama.CompressionSnappy, // Use compression for better performance
			}

			producer, err := NewBatchProducer(brokers, config)
			if err != nil {
				b.Fatalf("Failed to create producer: %v", err)
			}
			producer.Start()
			defer producer.Close()

			b.ResetTimer()
			start := time.Now()

			// Send messages
			for i := 0; i < numMessages; i++ {
				key := fmt.Sprintf("key-%d", i)
				err := producer.SendMessage(testTopic, key, message)
				if err != nil {
					b.Fatalf("Failed to send message: %v", err)
				}
			}

			//// Wait for all messages to be processed
			//producer.Close()
			//<-producer.Done()

			duration := time.Since(start)

			// Calculate and report metrics
			messagesPerSecond := float64(numMessages) / duration.Seconds()
			bytesPerSecond := float64(messageSize*numMessages) / duration.Seconds()

			b.ReportMetric(messagesPerSecond, "msgs/sec")
			b.ReportMetric(bytesPerSecond, "bytes/sec")

			// Report additional batch-specific metrics
			//msgSent, batchSent, errors := producer.Stats()
			//b.ReportMetric(float64(batchSent), "batches_sent")
			//b.ReportMetric(float64(errors), "errors")
			//b.ReportMetric(float64(msgSent)/float64(batchSent), "avg_batch_size")
		})
	}
}
