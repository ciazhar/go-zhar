package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

// Create a reusable async Kafka producer
func CreateAsyncProducer(brokers []string) sarama.AsyncProducer {
	config := sarama.NewConfig()

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

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to create async producer", nil)
	}

	return producer
}

// Send message asynchronously
func SendAsyncMessage(producer sarama.AsyncProducer, topic, message string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	// Send the message asynchronously
	producer.Input() <- msg
}

// Send message asynchronously
func SendAsyncMessageWithKey(producer sarama.AsyncProducer, topic, key, message string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}

	// Send the message asynchronously
	producer.Input() <- msg
}

func HandleAsyncResponse(producer sarama.AsyncProducer) {
	// Handle async successes and errors in separate goroutines
	go func() {
		for success := range producer.Successes() {
			logger.LogDebug(context.Background(), "Async message sent", map[string]interface{}{
				"partition": success.Partition,
				"offset":    success.Offset,
			})
		}
	}()

	go func() {
		for err := range producer.Errors() {
			logger.LogDebug(context.Background(), "Async message failed", map[string]interface{}{
				"partition": err.Msg.Partition,
				"offset":    err.Msg.Offset,
				"error":     err.Err.Error(),
			})
		}
	}()
}
