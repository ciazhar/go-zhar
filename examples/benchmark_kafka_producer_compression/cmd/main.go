package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/golang/snappy"
	"github.com/pierrec/lz4"
	"github.com/klauspost/compress/zstd"
)

type CompressionStats struct {
	OriginalSize     int
	CompressedSize   int
	CompressionTime  time.Duration
	CompressionRatio float64
}

func main() {
	// Test message
	message := generateTestMessage()
	topic := "compression-test"

	// Test different compression types
	compressionTypes := map[string]sarama.CompressionCodec{
		"none":   sarama.CompressionNone,
		"gzip":   sarama.CompressionGZIP,
		"snappy": sarama.CompressionSnappy,
		"lz4":    sarama.CompressionLZ4,
		"zstd":   sarama.CompressionZSTD,
	}

	results := make(map[string]CompressionStats)

	for name, codec := range compressionTypes {
		stats, err := measureCompression(topic, message, codec)
		if err != nil {
			log.Printf("Error measuring compression for %s: %v", name, err)
			continue
		}
		results[name] = stats
	}

	// Print results
	fmt.Println("\nCompression Results:")
	fmt.Println("-------------------")
	for name, stats := range results {
		fmt.Printf("Compression Type: %s\n", name)
		fmt.Printf("Original Size: %d bytes\n", stats.OriginalSize)
		fmt.Printf("Compressed Size: %d bytes\n", stats.CompressedSize)
		fmt.Printf("Compression Ratio: %.2f%%\n", stats.CompressionRatio)
		fmt.Printf("Compression Time: %v\n", stats.CompressionTime)
		fmt.Println("-------------------")
	}
}

func compressData(data []byte, codec sarama.CompressionCodec) ([]byte, error) {
	var compressed []byte
	var err error

	switch codec {
	case sarama.CompressionNone:
		return data, nil

	case sarama.CompressionGZIP:
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		if _, err := writer.Write(data); err != nil {
			return nil, err
		}
		if err := writer.Close(); err != nil {
			return nil, err
		}
		compressed = buf.Bytes()

	case sarama.CompressionSnappy:
		compressed = snappy.Encode(nil, data)

	case sarama.CompressionLZ4:
		var buf bytes.Buffer
		writer := lz4.NewWriter(&buf)
		if _, err := writer.Write(data); err != nil {
			return nil, err
		}
		if err := writer.Close(); err != nil {
			return nil, err
		}
		compressed = buf.Bytes()

	case sarama.CompressionZSTD:
		encoder, err := zstd.NewWriter(nil)
		if err != nil {
			return nil, err
		}
		compressed = encoder.EncodeAll(data, make([]byte, 0, len(data)))
		encoder.Close()

	default:
		return nil, fmt.Errorf("unknown compression codec: %v", codec)
	}

	return compressed, err
}

func measureCompression(topic, message string, codec sarama.CompressionCodec) (CompressionStats, error) {
	// Create producer config
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Compression = codec
	config.Producer.MaxMessageBytes = 20 * 1024 * 1024 // 20MB max message size

	// Measure original size
	messageBytes := []byte(message)
	originalSize := len(messageBytes)

	// Measure compression time and size
	start := time.Now()
	compressed, err := compressData(messageBytes, codec)
	if err != nil {
		return CompressionStats{}, fmt.Errorf("compression failed: %v", err)
	}
	compressionTime := time.Since(start)
	compressedSize := len(compressed)

	// Create producer
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		return CompressionStats{}, fmt.Errorf("failed to create producer: %v", err)
	}
	defer producer.Close()

	// Send message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return CompressionStats{}, fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)

	// Calculate compression ratio
	compressionRatio := float64(originalSize-compressedSize) / float64(originalSize) * 100

	return CompressionStats{
		OriginalSize:     originalSize,
		CompressedSize:   compressedSize,
		CompressionTime:  compressionTime,
		CompressionRatio: compressionRatio,
	}, nil
}

func generateTestMessage() string {
	// Generate a sample message with repeating content
	sampleText := "This is a test message with some repeating content. "
	var builder strings.Builder
	for i := 0; i < 1000; i++ {
		builder.WriteString(sampleText)
	}
	return builder.String()
}