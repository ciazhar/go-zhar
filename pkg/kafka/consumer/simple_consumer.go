package consumer

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"log"
	"strings"
	"sync"
	"time"
)

// SimpleConsumer optimized for parallel processing and retries.
// Pros: Simple to understand and use
// Cons: No automatic partition balancing
type SimpleConsumer struct {
	Topic      string
	Process    MessageProcessor // Function to process messages
	Partition  int32
	Offset     int64
	consumer   sarama.Consumer
	maxRetries int
	workers    int
}

// NewSimpleConsumer creates a new SimpleConsumer instance.
func NewSimpleConsumer(brokerList []string, topic string, process MessageProcessor, assignor string, offsetOldest bool, workers int) *SimpleConsumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Retry.Backoff = 2 * time.Second

	if offsetOldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

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

	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to create Kafka consumer", nil)
		return nil
	}

	return &SimpleConsumer{
		consumer:   consumer,
		maxRetries: 3,
		Topic:      topic,
		Process:    process,
		workers:    workers,
	}
}

// ConsumeWithRetry starts consuming messages with retry logic and parallel processing.
func (c *SimpleConsumer) ConsumeWithRetry(topic string, partition int32, offset int64) error {
	var retries int
	var partitionConsumer sarama.PartitionConsumer
	var err error

	for retries < c.maxRetries {
		partitionConsumer, err = c.consumer.ConsumePartition(topic, partition, offset)
		if err == nil {
			break
		}
		retries++
		time.Sleep(time.Second * time.Duration(1<<retries)) // Exponential backoff
	}

	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	// Worker pool for parallel processing
	messageCh := make(chan *sarama.ConsumerMessage, 100)
	var wg sync.WaitGroup

	for i := 0; i < c.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for message := range messageCh {
				if err := c.Process(message); err != nil {
					log.Printf("Error processing message: %v\n", err)
				}
			}
		}()
	}

	for message := range partitionConsumer.Messages() {
		messageCh <- message
	}
	close(messageCh)
	wg.Wait() // Wait for all workers to complete

	return nil
}

// StartSimpleConsumer starts consuming messages from Kafka.
func StartSimpleConsumer(
	ctx context.Context,
	brokers string,
	consumers map[string]SimpleConsumer,
	wg *sync.WaitGroup,
	assignor string,
	offsetOldest bool,
	workers int,
) {

	for topic, config := range consumers {
		kafkaConsumer := NewSimpleConsumer(strings.Split(brokers, ","), config.Topic, config.Process, assignor, offsetOldest, workers)
		wg.Add(1)
		go func() {
			err := kafkaConsumer.ConsumeWithRetry(topic, config.Partition, config.Offset)
			if err != nil {
				logger.LogFatal(ctx, err, "Failed to consume messages", nil)
			}
		}()
		logger.LogInfo(ctx, "Started consumer for topic", map[string]interface{}{"topic": topic})
	}
}
