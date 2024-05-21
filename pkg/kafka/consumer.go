package kafka

import (
	"github.com/IBM/sarama"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"os"
	"os/signal"
	"strings"
)

type Consumer struct {
	logger   *logger.Logger
	consumer sarama.Consumer
}

func NewConsumer(brokers string, logger *logger.Logger) Consumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(strings.Split(brokers, ","), config)
	if err != nil {
		logger.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	logger.Info("Connected to Kafka Consumer")
	return Consumer{
		logger:   logger,
		consumer: consumer,
	}
}

func (c *Consumer) ConsumeMessages(topicName string, out func(msg string)) {

	consumer, err := c.consumer.ConsumePartition(topicName, 0, sarama.OffsetOldest)
	if err != nil {
		c.logger.Fatalf("Error creating partition consumer: %v", err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			c.logger.Fatalf("Error closing partition consumer: %v", err)
		}
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				c.logger.Infof("Error from consumer: %v", err)
			case msg := <-consumer.Messages():
				out(string(msg.Value))
			case <-signals:
				c.logger.Info("Interrupt signal detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
}

func (c *Consumer) Close() {
	err := c.consumer.Close()
	if err != nil {
		c.logger.Fatalf("Failed to close Kafka consumer: %v", err)
	}
}
