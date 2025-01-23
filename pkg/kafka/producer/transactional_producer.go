package producer

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"sync"
	"sync/atomic"
	"time"
)

// TransactionalProducer TRANSACTIONAL PRODUCER
// Pros: Atomic multi-partition writes
// Cons: Higher latency, more complex
type TransactionalProducer struct {
	producer      sarama.AsyncProducer
	messagesSent  uint64
	batchesSent   uint64
	errors        uint64
	mu            sync.RWMutex
	closed        bool
	txnInProgress bool
	wg            sync.WaitGroup
}

// Message represents a generic Kafka message
type Message struct {
	Topic   string      `json:"topic"`
	Key     string      `json:"key"`
	Value   interface{} `json:"value"`
	Headers []Header    `json:"headers,omitempty"`
}

// Header represents a Kafka message header
type Header struct {
	Key   string
	Value []byte
}

// MessageBatch represents a batch of messages to be sent in a transaction
type MessageBatch struct {
	Messages []Message
	Options  *BatchOptions
}

// BatchOptions contains additional options for batch processing
type BatchOptions struct {
	TransactionID string
	Timeout       time.Duration
	Retries       int
}

func NewTransactionalProducer(brokers []string, transactionalID string) (*TransactionalProducer, error) {
	config := sarama.NewConfig()

	// Transaction settings
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Idempotent = true
	config.Producer.Transaction.ID = transactionalID
	config.Producer.Transaction.Retry.Max = 3
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Net.MaxOpenRequests = 1

	// Performance optimizations
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.CompressionLevel = 1
	config.Producer.MaxMessageBytes = 1000000
	config.Producer.Flush.Bytes = 64 * 1024
	config.Producer.Flush.Messages = 100
	config.Producer.Flush.Frequency = 1 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	tp := &TransactionalProducer{
		producer: producer,
		closed:   false,
	}

	tp.handleAsyncResults()
	return tp, nil
}

func (p *TransactionalProducer) handleAsyncResults() {
	p.wg.Add(2)

	// Success handler
	go func() {
		defer p.wg.Done()
		for range p.producer.Successes() {
			atomic.AddUint64(&p.messagesSent, 1)
		}
	}()

	// Error handler
	go func() {
		defer p.wg.Done()
		for err := range p.producer.Errors() {
			if err != nil {
				atomic.AddUint64(&p.errors, 1)
			}
		}
	}()
}

// SendMessage sends a single message
func (p *TransactionalProducer) SendMessage(msg Message) error {
	return p.SendMessageBatch(MessageBatch{
		Messages: []Message{msg},
	})
}

// SendMessageBatch sends multiple messages in a transaction
func (p *TransactionalProducer) SendMessageBatch(batch MessageBatch) error {
	if len(batch.Messages) == 0 {
		return ErrEmptyBatch
	}

	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return ErrProducerClosed
	}
	if p.txnInProgress {
		p.mu.Unlock()
		return ErrTxnInProgress
	}
	p.txnInProgress = true
	p.mu.Unlock()

	defer func() {
		p.mu.Lock()
		p.txnInProgress = false
		p.mu.Unlock()
	}()

	// Begin transaction
	if err := p.producer.BeginTxn(); err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Process messages
	success := true
	var errs []error

	for _, msg := range batch.Messages {
		producerMsg, err := p.createProducerMessage(msg)
		if err != nil {
			success = false
			errs = append(errs, err)
			break
		}

		p.producer.Input() <- producerMsg
	}

	// Commit or abort transaction
	if success {
		if err := p.producer.CommitTxn(); err != nil {
			abortErr := p.producer.AbortTxn()
			if abortErr != nil {
				return fmt.Errorf("failed to commit and abort transaction: commit err: %w, abort err: %v", err, abortErr)
			}
			return fmt.Errorf("failed to commit transaction (aborted): %w", err)
		}
		atomic.AddUint64(&p.messagesSent, uint64(len(batch.Messages)))
		atomic.AddUint64(&p.batchesSent, 1)
	} else {
		if err := p.producer.AbortTxn(); err != nil {
			return fmt.Errorf("failed to abort transaction: %w", err)
		}
		return fmt.Errorf("transaction aborted due to errors: %v", errs)
	}

	return nil
}

// createProducerMessage converts a Message to sarama.ProducerMessage
func (p *TransactionalProducer) createProducerMessage(msg Message) (*sarama.ProducerMessage, error) {
	if msg.Topic == "" {
		return nil, ErrInvalidMessage
	}

	valueBytes, err := json.Marshal(msg.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message value: %w", err)
	}

	producerMsg := &sarama.ProducerMessage{
		Topic: msg.Topic,
		Value: sarama.ByteEncoder(valueBytes),
	}

	if msg.Key != "" {
		producerMsg.Key = sarama.StringEncoder(msg.Key)
	}

	if len(msg.Headers) > 0 {
		producerMsg.Headers = make([]sarama.RecordHeader, len(msg.Headers))
		for i, h := range msg.Headers {
			producerMsg.Headers[i] = sarama.RecordHeader{
				Key:   []byte(h.Key),
				Value: h.Value,
			}
		}
	}

	return producerMsg, nil
}

func (p *TransactionalProducer) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	p.mu.Unlock()

	// Close the producer
	if err := p.producer.Close(); err != nil {
		return err
	}

	// Wait for all messages to be processed
	p.wg.Wait()
	return nil
}
