package consumer

import (
	"github.com/IBM/sarama"
	"sync"
)

// CompactedTopicConsumer COMPACTED TOPIC CONSUMER
// Pros: Efficient for key-value storage, reduced storage usage
// Cons: Only keeps latest value per key
type CompactedTopicConsumer struct {
	consumer sarama.Consumer
	store    map[string][]byte
	mutex    sync.RWMutex
}

func NewCompactedTopicConsumer(brokerList []string) (*CompactedTopicConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	// Set fetch size to optimize for compacted topics
	config.Consumer.Fetch.Max = 1024 * 1024 // 1MB

	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return &CompactedTopicConsumer{
		consumer: consumer,
		store:    make(map[string][]byte),
	}, nil
}

func (c *CompactedTopicConsumer) ConsumeCompactedTopic(topic string) error {
	partitions, err := c.consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for _, partition := range partitions {
		pc, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				c.mutex.Lock()
				if msg.Value == nil {
					// Tombstone message, delete the key
					delete(c.store, string(msg.Key))
				} else {
					// Update the value for the key
					c.store[string(msg.Key)] = msg.Value
				}
				c.mutex.Unlock()
			}
		}(pc)
	}

	return nil
}
