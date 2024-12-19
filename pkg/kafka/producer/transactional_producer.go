package producer

import (
	"fmt"
	"github.com/IBM/sarama"
)

// TransactionalProducer TRANSACTIONAL PRODUCER
// Pros: Atomic multi-partition writes
// Cons: Higher latency, more complex
type TransactionalProducer struct {
	producer sarama.AsyncProducer
}

func NewTransactionalProducer(brokerList []string, transactionalID string) (*TransactionalProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Idempotent = true
	config.Producer.Transaction.ID = transactionalID
	config.Producer.Transaction.Retry.Max = 3
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Net.MaxOpenRequests = 1 // Required for idempotence

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return &TransactionalProducer{producer: producer}, nil
}

func (p *TransactionalProducer) SendMessagesInTransaction(topic string, messages []string) error {
	// Begin transaction
	err := p.producer.BeginTxn()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Send messages
	for _, msg := range messages {
		p.producer.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(msg),
		}
	}

	// Commit transaction
	err = p.producer.CommitTxn()
	if err != nil {
		// If commit fails, attempt to abort the transaction
		abortErr := p.producer.AbortTxn()
		if abortErr != nil {
			return fmt.Errorf("failed to commit transaction: %v and failed to abort: %v", err, abortErr)
		}
		return fmt.Errorf("failed to commit transaction (successfully aborted): %v", err)
	}

	return nil
}

func (p *TransactionalProducer) Close() error {
	return p.producer.Close()
}
