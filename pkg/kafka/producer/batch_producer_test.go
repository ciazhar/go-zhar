package producer

import (
	"fmt"
	"github.com/IBM/sarama"
	"testing"
)

func BenchmarkBatchProducer_SendMessage(b *testing.B) {
	brokers := []string{"localhost:9092"}
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		producer.Input() <- Message{
			Topic: "benchmark_topic",
			Key:   []byte("key"),
			Value: []byte("value"),
		}
	}
}

func BenchmarkBatchProducer_BatchSizes(b *testing.B) {
	batchSizes := []int{1, 10, 100, 1000}
	brokers := []string{"localhost:9092"}

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
				producer.Input() <- Message{
					Topic: "benchmark_topic",
					Key:   []byte("key"),
					Value: []byte("value"),
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
	brokers := []string{"localhost:9092"}

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

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				producer.Input() <- Message{
					Topic: "benchmark_topic",
					Key:   []byte("key"),
					Value: []byte("value"),
				}
			}
		})
	}
}
