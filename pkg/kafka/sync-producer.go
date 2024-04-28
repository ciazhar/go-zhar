package kafka

import (
	"github.com/IBM/sarama"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"strings"
)

type SyncProducer struct {
	producer sarama.SyncProducer
	logger   logger.Logger
}

type SyncProducerConfig struct {
	Version string
	Retry   int
}

func NewSyncProducer(brokers string, logger logger.Logger, producerConfig ...SyncProducerConfig) *SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	config.Producer.MaxMessageBytes = 33554432

	if len(producerConfig) > 0 {
		version, err := sarama.ParseKafkaVersion(producerConfig[0].Version)
		if err != nil {
			logger.Fatalf("Error parsing Kafka version: %v", err)
		}
		config.Version = version
		config.Producer.Retry.Max = producerConfig[0].Retry
	}

	producer, err := sarama.NewSyncProducer(strings.Split(brokers, ","), config)
	if err != nil {
		logger.Fatalf("Failed to create Kafka Producer: %v", err)
	}
	logger.Info("Connected to Kafka Producer")
	return &SyncProducer{
		producer: producer,
		logger:   logger,
	}
}

func (p *SyncProducer) PublishMessage(topic string, value string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(value),
	}

	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		p.logger.Infof("Failed to publish a message: %v", err)
	}
}

func (p *SyncProducer) PublishMessageWithKey(topic string, key string, message string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}

	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		p.logger.Infof("Failed to publish a message: %v", err)
	}
}

func (p *SyncProducer) Close() {
	err := p.producer.Close()
	if err != nil {
		p.logger.Fatalf("Failed to close Kafka Producer: %v", err)
	}
}
