package kafka

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

// MessageProcessor is a function type for processing messages.
type MessageProcessor func(msg *sarama.ConsumerMessage) error

// KafkaConsumer wraps a Sarama consumer group and a message processor.
type KafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	topic         string
	process       MessageProcessor // Function to process messages
}

// NewKafkaConsumer creates a new KafkaConsumer.
func NewKafkaConsumer(brokers []string, groupID, topic string, process MessageProcessor) (*KafkaConsumer, error) {
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, nil)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to create Kafka consumer group", nil)
	}

	return &KafkaConsumer{
		consumerGroup: consumerGroup,
		topic:         topic,
		process:       process, // Assign the processing function
	}, nil
}

// ConsumeMessages starts consuming messages from Kafka.
func (kc *KafkaConsumer) ConsumeMessages(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			logger.LogInfo(context.Background(), "Shutting down Kafka consumer group...", nil)
			return
		default:
			if err := kc.consumerGroup.Consume(ctx, []string{kc.topic}, kc); err != nil {
				logger.LogWarn(context.Background(), err, "Failed to consume messages from Kafka", nil)
				return
			}
		}
	}
}

// Setup is run at the beginning of a new session, not used here.
func (kc *KafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, not used here.
func (kc *KafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from a single partition.
func (kc *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// Call the processing function
		if err := kc.process(msg); err != nil {
			logger.LogWarn(context.Background(), err, "Failed to process message", nil)
		}
		session.MarkMessage(msg, "") // Mark message as processed
	}
	return nil
}

// Close closes the Kafka consumer group.
func (kc *KafkaConsumer) Close() error {
	logger.LogInfo(context.Background(), "Closing Kafka consumer group...", nil)
	return kc.consumerGroup.Close()
}
