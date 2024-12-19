package consumer

import (
	"github.com/IBM/sarama"
	"log"
	"time"
)

// SimpleConsumer SIMPLE CONSUMER WITH RETRY
// Pros: Simple to understand and use
// Cons: No automatic partition balancing
type SimpleConsumer struct {
	consumer   sarama.Consumer
	maxRetries int
}

func NewSimpleConsumer(brokerList []string) (*SimpleConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Retry.Backoff = 2 * time.Second

	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return &SimpleConsumer{
		consumer:   consumer,
		maxRetries: 3,
	}, nil
}

func (c *SimpleConsumer) ConsumeWithRetry(topic string, partition int32) error {
	var retries int
	var partitionConsumer sarama.PartitionConsumer
	var err error

	for retries < c.maxRetries {
		partitionConsumer, err = c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err == nil {
			break
		}
		retries++
		time.Sleep(time.Second * time.Duration(retries))
	}

	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		log.Printf("Message topic:%q partition:%d offset:%d\n",
			message.Topic, message.Partition, message.Offset)
	}

	return nil
}
