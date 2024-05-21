package kafka

import (
	"errors"
	"github.com/IBM/sarama"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"strings"
)

type Admin struct {
	logger *logger.Logger
	admin  sarama.ClusterAdmin
}

func NewAdmin(brokers string, logger *logger.Logger) Admin {
	config := sarama.NewConfig()

	admin, err := sarama.NewClusterAdmin(strings.Split(brokers, ","), config)
	if err != nil {
		logger.Fatalf("Failed to create Kafka admin client: %v", err)
	}

	return Admin{
		logger: logger,
		admin:  admin,
	}
}

type CreateTopicConfig struct {
	NumPartitions     int32
	ReplicationFactor int16
}

func (a *Admin) CreateTopic(topicName string, config ...CreateTopicConfig) {
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     1,   // Number of partitions for the topic
		ReplicationFactor: 1,   // Replication factor for the topic
		ConfigEntries:     nil, // Additional topic configuration (can be nil)
	}

	if len(config) > 0 {
		topicDetail.NumPartitions = config[0].NumPartitions
		topicDetail.ReplicationFactor = config[0].ReplicationFactor
	}

	// Create the topic
	err := a.admin.CreateTopic(topicName, topicDetail, true)
	if err != nil {
		if errors.Is(err, sarama.ErrTopicAlreadyExists) {
			a.logger.Infof("Topic '%s' already exists.", topicName)
			return
		} else {
			a.logger.Fatalf("Failed to create Kafka topic: %v", err)
		}
	}

	a.logger.Infof("Topic '%s' created successfully.", topicName)
}

func (c *Admin) Close() {
	err := c.admin.Close()
	if err != nil {
		c.logger.Fatalf("Failed to close Kafka admin client: %v", err)
	}
}
