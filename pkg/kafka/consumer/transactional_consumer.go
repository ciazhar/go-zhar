package consumer

import (
	"context"
	"github.com/IBM/sarama"
)

// TransactionalConsumer TRANSACTIONAL CONSUMER
// Pros: Exactly-once processing semantics
// Cons: Higher latency, more complex setup
type TransactionalConsumer struct {
	consumer sarama.ConsumerGroup
	producer sarama.AsyncProducer // For sending results transactionally
}

func NewTransactionalConsumer(brokerList []string, groupID string, transactionalID string) (*TransactionalConsumer, error) {
	// Consumer config
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	consumerConfig.Consumer.IsolationLevel = sarama.ReadCommitted // Only read committed messages

	// Producer config for transactional output
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Transaction.ID = transactionalID
	producerConfig.Producer.Idempotent = true
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll

	consumer, err := sarama.NewConsumerGroup(brokerList, groupID, consumerConfig)
	if err != nil {
		return nil, err
	}

	producer, err := sarama.NewAsyncProducer(brokerList, producerConfig)
	if err != nil {
		consumer.Close()
		return nil, err
	}

	return &TransactionalConsumer{
		consumer: consumer,
		producer: producer,
	}, nil
}

func (c *TransactionalConsumer) ConsumeAndProcess(ctx context.Context, topics []string, processor func(msg *sarama.ConsumerMessage) (*sarama.ProducerMessage, error)) error {
	handler := &TransactionalConsumerHandler{
		producer:  c.producer,
		processor: processor,
	}

	for {
		err := c.consumer.Consume(ctx, topics, handler)
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

type TransactionalConsumerHandler struct {
	producer  sarama.AsyncProducer
	processor func(msg *sarama.ConsumerMessage) (*sarama.ProducerMessage, error)
}

func (h *TransactionalConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *TransactionalConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *TransactionalConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := h.producer.BeginTxn(); err != nil {
			return err
		}

		outputMsg, err := h.processor(msg)
		if err != nil {
			h.producer.AbortTxn()
			continue
		}

		h.producer.Input() <- outputMsg

		if err := h.producer.CommitTxn(); err != nil {
			return err
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
