package kafka

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type ConsumerGroup struct {
	logger        logger.Logger
	consumerGroup sarama.ConsumerGroup
	ready         chan bool
}

type ConsumerGroupConfig struct {
	Version       string
	InitialOffset string
	Assignor      string
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *ConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *ConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (c *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func NewConsumerGroup(brokers string, groupId string, logger logger.Logger, conf ...ConsumerGroupConfig) ConsumerGroup {
	config := sarama.NewConfig()
	config.Consumer.Offsets.AutoCommit.Enable = false // disable auto-commit

	if len(conf) > 0 {
		version, err := sarama.ParseKafkaVersion(conf[0].Version)
		if err != nil {
			logger.Fatalf("Error parsing Kafka version: %v", err)
		}

		config.Version = version
		assignor := strings.ToLower(conf[0].Assignor)
		switch assignor {
		case "sticky":
			config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
		case "roundrobin":
			config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
		case "range":
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

func (c *ConsumerGroup) ConsumeMessages(topic string, out func(msg string)) {
	keepRunning := true
	c.ready = make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())
	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := c.consumerGroup.Consume(ctx, strings.Split(topic, ","), c); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				c.logger.Fatalf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			c.ready = make(chan bool)
		}
	}()

	<-c.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			c.logger.Info("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			c.logger.Info("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			c.toggleConsumptionFlow(c.consumerGroup, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err := c.consumerGroup.Close(); err != nil {
		c.logger.Fatalf("Error closing client: %v", err)
	}
}

func (c *ConsumerGroup) toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		c.logger.Info("Resuming consumption")
	} else {
		client.PauseAll()
		c.logger.Info("Pausing consumption")
	}

	*isPaused = !*isPaused
}
