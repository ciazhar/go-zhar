package kafka

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"log"
	"strings"
	"sync"
)

type ConsumerConfig struct {
	Topics  []string
	Handler sarama.ConsumerGroupHandler
}

func StartConsumers(ctx context.Context, brokers string, consumers map[string]ConsumerConfig, logger logger.Logger) {
	// Create a wait group to synchronize the consumer goroutines
	var wg sync.WaitGroup

	// Iterate over the consumers map
	for groupID, config := range consumers {

		consumerGroup := NewConsumerGroup(brokers, groupID, logger)

		wg.Add(1)

		go runConsumer(ctx, &wg, consumerGroup, config.Topics, config.Handler)

		logger.Infof("Starting consumer group %s topic %s", groupID, config.Topics)
	}

	// Wait for all the consumer goroutines to finish
	wg.Wait()
}

func runConsumer(
	ctx context.Context,
	wg *sync.WaitGroup,
	consumerGroup ConsumerGroup,
	topics []string,
	handler sarama.ConsumerGroupHandler,
) {
	defer wg.Done()
	defer consumerGroup.Close()

	for {
		err := consumerGroup.GetInstance().Consume(ctx, topics, handler)
		if err != nil {
			log.Println("Error consuming messages:", err)
		}

		if errors.Is(err, sarama.ErrClosedConsumerGroup) {
			break
		}
	}
}

type ConsumerGroup struct {
	logger        logger.Logger
	consumerGroup sarama.ConsumerGroup
	ready         chan bool
	out           func(key, value string)
}

type ConsumerGroupConfig struct {
	Version       string
	InitialOffset string
	Assignor      string
}

func NewConsumerGroup(brokers string, groupId string, logger logger.Logger, conf ...ConsumerGroupConfig) ConsumerGroup {
	config := sarama.NewConfig()
	config.Consumer.Offsets.AutoCommit.Enable = false // disable auto-commit
	config.Consumer.Return.Errors = true

	if len(conf) > 0 {
		if conf[0].Version != "" {
			version, err := sarama.ParseKafkaVersion(conf[0].Version)
			if err != nil {
				logger.Fatalf("Error parsing Kafka version: %v", err)
			}

			config.Version = version
		}

		assignor := strings.ToLower(conf[0].Assignor)
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

		if conf[0].InitialOffset == "oldest" {
			config.Consumer.Offsets.Initial = sarama.OffsetOldest
		}
	}

	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), groupId, config)
	if err != nil {
		logger.Fatalf("Error creating consumer group client: %v", err)
	}
	return ConsumerGroup{
		consumerGroup: client,
	}
}

func (c *ConsumerGroup) GetInstance() sarama.ConsumerGroup {
	return c.consumerGroup
}

func (c *ConsumerGroup) Close() {
	err := c.consumerGroup.Close()
	if err != nil {
		c.logger.Fatalf("Failed to close Kafka consumer group: %v", err)
	}
}

type BasicConsumerHandler struct {
	Ready chan bool
	Out   func(key, value string)
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *BasicConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *BasicConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (c *BasicConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}

			c.Out(string(message.Key), string(message.Value))

			session.MarkMessage(message, "")

			session.Commit()

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func NewBasicConsumerHandler(out func(key, value string)) *BasicConsumerHandler {
	return &BasicConsumerHandler{
		Ready: make(chan bool),
		Out:   out,
	}
}
