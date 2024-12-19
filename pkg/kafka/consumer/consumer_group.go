package consumer

import (
	"context"
	"github.com/IBM/sarama"
)

// ConsumerGroup CONSUMER GROUP
// Pros: Automatic partition balancing, parallel processing
// Cons: More complex than simple consumer
type ConsumerGroup struct {
	group   sarama.ConsumerGroup
	topics  []string
	handler sarama.ConsumerGroupHandler
}

type ConsumerGroupHandler struct {
	ready chan bool
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		session.MarkMessage(message, "")
	}
	return nil
}

func NewConsumerGroup(brokerList []string, groupID string, topics []string) (*ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup(brokerList, groupID, config)
	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{
		group:   group,
		topics:  topics,
		handler: &ConsumerGroupHandler{ready: make(chan bool)},
	}, nil
}

func (c *ConsumerGroup) Consume(ctx context.Context) error {
	for {
		err := c.group.Consume(ctx, c.topics, c.handler)
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
