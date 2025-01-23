package producer

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"sync"
	"time"
)

// IdempotentProducer IDEMPOTENT PRODUCER
// Pros: Exactly-once delivery semantics
// Cons: Higher latency, requires newer Kafka version
type IdempotentProducer struct {
	producer sarama.SyncProducer
	mu       sync.RWMutex
	closed   bool
}

func NewIdempotentProducer(brokerList []string) (*IdempotentProducer, error) {
	config := sarama.NewConfig()

	// Idempotent producer configuration
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Idempotent = true
	config.Producer.Return.Successes = true
	config.Net.MaxOpenRequests = 1

	// Performance optimizations
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.CompressionLevel = 1
	config.Producer.MaxMessageBytes = 1000000
	config.Producer.Retry.Max = 3
	config.Producer.Retry.Backoff = 100 * time.Millisecond

	//// Batch settings
	//config.Producer.Flush.Bytes = 64 * 1024
	//config.Producer.Flush.Messages = 100
	//config.Producer.Flush.Frequency = 1 * time.Millisecond

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return &IdempotentProducer{
		producer: producer,
		closed:   false,
	}, nil
}

func (p *IdempotentProducer) SendMessage(topic, key string, value any) (partition int32, offset int64, err error) {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return 0, 0, ErrProducerClosed
	}
	p.mu.RUnlock()

	marshal, err := json.Marshal(value)
	if err != nil {
		return 0, 0, err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(marshal),
	}

	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}

	partition, offset, err = p.producer.SendMessage(msg)
	if err != nil {
		return partition, offset, err
	}

	return partition, offset, nil
}

func (p *IdempotentProducer) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	p.mu.Unlock()

	return p.producer.Close()
}
