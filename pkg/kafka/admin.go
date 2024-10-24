package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

// CreateKafkaAdminClient creates and returns a single Kafka admin client.
func CreateKafkaAdminClient(brokers []string) (sarama.ClusterAdmin, error) {
	configSarama := sarama.NewConfig()

	admin, err := sarama.NewClusterAdmin(brokers, configSarama)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to create Kafka admin client", nil)
	}
	return admin, nil
}

// CreateKafkaTopic creates a new Kafka topic using an existing admin client.
func CreateKafkaTopic(admin sarama.ClusterAdmin, topicName string, numPartitions int32, replicationFactor int16, retentionMs int64, config map[string]string) error {
	// Prepare topic details
	topicDetail := sarama.TopicDetail{
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
		ConfigEntries:     make(map[string]*string),
	}

	// Set retention policy and other configurations
	topicDetail.ConfigEntries["retention.ms"] = strPtr(fmt.Sprintf("%d", retentionMs))
	for key, value := range config {
		topicDetail.ConfigEntries[key] = strPtr(value)
	}

	// Create the topic
	err := admin.CreateTopic(topicName, &topicDetail, false)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to create topic", map[string]interface{}{
			"topic": topicName,
		})
	}

	return nil
}

// Helper function to convert string to *string
func strPtr(s string) *string {
	return &s
}