package producer

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"sync"
	"time"
)

// SyncProducer SYNCHRONOUS PRODUCER
// Pros: Simple to use, immediate feedback on success/failure
// Cons: Lower throughput due to waiting for response
type SyncProducer struct {
	producer sarama.SyncProducer
	mu       sync.RWMutex
	closed   bool
}

func NewSyncProducer(brokerList []string, maxRetry int) (*SyncProducer, error) {
	config := sarama.NewConfig()

	// Basic configuration
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = maxRetry

	// Performance optimizations
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.CompressionLevel = 1
	config.Producer.MaxMessageBytes = 1000000
	config.Producer.Retry.Backoff = 100 * time.Millisecond

	// Batch settings for better throughput
	//config.Producer.Flush.Bytes = 64 * 1024
	//config.Producer.Flush.Messages = 100
	//config.Producer.Flush.Frequency = 1 * time.Millisecond

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return &SyncProducer{
		producer: producer,
		closed:   false,
	}, nil
}

func (p *SyncProducer) SendMessage(topic, key string, value any) (partition int32, offset int64, err error) {
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

func (p *SyncProducer) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	p.mu.Unlock()

	return p.producer.Close()
}
