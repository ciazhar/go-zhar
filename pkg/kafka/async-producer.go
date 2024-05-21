package kafka

import (
	"github.com/IBM/sarama"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"strings"
	"time"
)

type AsyncProducer struct {
	logger   *logger.Logger
	producer sarama.AsyncProducer
}

type AsyncProducerConfig struct {
	Version string
	Retry   int
}

func NewAsyncProducer(brokers string, logger *logger.Logger, producerConfig ...AsyncProducerConfig) *AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	if len(producerConfig) > 0 {
		version, err := sarama.ParseKafkaVersion(producerConfig[0].Version)
		if err != nil {
			logger.Fatalf("Error parsing Kafka version: %v", err)
		}
		config.Version = version
		config.Producer.Retry.Max = producerConfig[0].Retry
	}

	producer, err := sarama.NewAsyncProducer(strings.Split(brokers, ","), config)
	if err != nil {
		logger.Fatalf("Failed to create Kafka producer: %v", err)
	}

	logger.Info("Connected to Kafka Producer")

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			logger.Infof("Failed to write access log entry: %v", err)
		}
	}()

	return &AsyncProducer{
		logger:   logger,
		producer: producer,
	}
}

func (p *AsyncProducer) PublishMessage(topic string, value string) {
	p.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(value),
	}
}

func (p *AsyncProducer) PublishMessageWithKey(topic string, key string, value string) {
	p.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(value),
		Key:   sarama.StringEncoder(key),
	}
}

func (p *AsyncProducer) Close() {
	err := p.producer.Close()
	if err != nil {
		p.logger.Fatalf("Failed to close Kafka producer: %v", err)
	}
}
