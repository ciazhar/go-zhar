package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

// Create a reusable Kafka producer
func CreateProducer(brokers []string) sarama.SyncProducer {
	config := sarama.NewConfig()

	// Best Practice: Enable return of successes
	config.Producer.Return.Successes = true

	// Best Practice: Enable idempotent producer to ensure exactly-once semantics, meaning even in case of retries, duplicates are avoided
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1

	// Best Practice: Set acknowledgments to all replicas to ensure durability
	config.Producer.RequiredAcks = sarama.WaitForAll

	// Best Practice: Set retries and backoff for transient failures
	config.Producer.Retry.Max = 5       // Retry up to 5 times
	config.Producer.Retry.Backoff = 100 // Wait 100ms between retries

	// Best Practice: Enable compression to reduce message size and improve performance
	config.Producer.Compression = sarama.CompressionGZIP

	// Create a sync producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to create Kafka producer", nil)
	}

	return producer
}

// Send message using the reusable producer
func SendMessage(producer sarama.SyncProducer, topic, message string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err := producer.SendMessage(msg)
	if err != nil {
		logger.LogAndReturnError(context.Background(), err, "Failed to send message", nil)
	}
}

func SendMessageWithKey(producer sarama.SyncProducer, topic, key, message string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}

	_, _, err := producer.SendMessage(msg)
	if err != nil {
		logger.LogAndReturnError(context.Background(), err, "Failed to send message", nil)
	}
}
