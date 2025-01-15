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
	//producer.Start()
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
			//producer.Start()
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
			//producer.Start()
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
	config := ProducerConfig{
		BatchSize:   1000,
		Compression: sarama.CompressionSnappy,
	}

	producer, err := NewBatchProducer(brokers, config)
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

func BenchmarkBatchProducer_DifferentSizes(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000} // 100B, 1KB, 10KB, 100KB
	config := ProducerConfig{
		BatchSize:   1000,
		Compression: sarama.CompressionSnappy,
	}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("MessageSize-%dB", size), func(b *testing.B) {
			brokers := []string{"localhost:9092", "localhost:9093", "localhost:9094"}
			producer, err := NewBatchProducer(brokers, config)
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
