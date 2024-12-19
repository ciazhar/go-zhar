package consumer

import (
	"github.com/IBM/sarama"
	"log"
	"sync"
)

// StreamConsumer KAFKA STREAMS CONSUMER
// Pros: Powerful stream processing capabilities
// Cons: More complex setup, higher resource usage
type StreamConsumer struct {
	consumer sarama.Consumer
	streams  map[string]chan *sarama.ConsumerMessage
	wg       sync.WaitGroup
}

func NewStreamConsumer(brokerList []string) (*StreamConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return &StreamConsumer{
		consumer: consumer,
		streams:  make(map[string]chan *sarama.ConsumerMessage),
	}, nil
}

func (c *StreamConsumer) CreateStream(topic string, processor func(*sarama.ConsumerMessage) error) error {
	partitions, err := c.consumer.Partitions(topic)
	if err != nil {
		return err
	}

	c.streams[topic] = make(chan *sarama.ConsumerMessage, 1000)

	for _, partition := range partitions {
		c.wg.Add(1)
		go func(p int32) {
			defer c.wg.Done()

			pc, err := c.consumer.ConsumePartition(topic, p, sarama.OffsetNewest)
			if err != nil {
				log.Printf("Failed to start consumer for partition %d: %s\n", p, err)
				return
			}
			defer pc.Close()

			for msg := range pc.Messages() {
				if err := processor(msg); err != nil {
					log.Printf("Error processing message: %v\n", err)
				}
			}
		}(partition)
	}

	return nil
}
