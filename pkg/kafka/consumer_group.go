package kafka

import (
	"context"
	"strings"
	"sync"

	"github.com/IBM/sarama"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

// MessageProcessor is a function type for processing messages.
type MessageProcessor func(msg *sarama.ConsumerMessage) error

// KafkaConsumer wraps a Sarama consumer group and a message processor.
type KafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	Topic         string
	GroupID       string
	Process       MessageProcessor // Function to process messages
}

// NewKafkaConsumer creates a new KafkaConsumer.
func NewKafkaConsumer(brokers []string, groupID, topic string, process MessageProcessor, assignor string, offsetOldest bool) *KafkaConsumer {

	config := sarama.NewConfig()
	config.Consumer.Offsets.AutoCommit.Enable = false // disable auto-commit
	config.Consumer.Return.Errors = true

	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	}

	if offsetOldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to create Kafka consumer group", nil)
	}

	return &KafkaConsumer{
		consumerGroup: consumerGroup,
		Topic:         topic,
		GroupID:       groupID,
		Process:       process, // Assign the processing function
	}
}

// ConsumeMessages starts consuming messages from Kafka.
func (kc *KafkaConsumer) ConsumeMessages(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer kc.Close()

	for {
		select {
		case <-ctx.Done():
			logger.LogInfo(context.Background(), "Shutting down Kafka consumer group...", nil)
			return
		default:
			if err := kc.consumerGroup.Consume(ctx, []string{kc.Topic}, kc); err != nil {
				logger.LogWarn(context.Background(), err, "Failed to consume messages from Kafka", nil)
				return
			}
		}
	}
}

// StartConsumers starts consuming messages from Kafka.
func StartConsumers(
	ctx context.Context,
	brokers string,
	consumers map[string]KafkaConsumer,
	wg *sync.WaitGroup,
	assignor string,
	offsetOldest bool,
) {

	for topic, config := range consumers {
		kafkaConsumer := NewKafkaConsumer(strings.Split(brokers, ","), config.GroupID, topic, config.Process, assignor, offsetOldest)
		wg.Add(1)
		go kafkaConsumer.ConsumeMessages(ctx, wg)
		logger.LogInfo(context.Background(), "Started consumer for topic", map[string]interface{}{"topic": topic, "group": config.GroupID})
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
		if err := kc.Process(msg); err != nil {
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
