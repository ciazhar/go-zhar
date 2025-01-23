package producer

import (
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

// BatchProducer BATCH PRODUCER
// Pros: Improved throughput, better network utilization
// Cons: Increased latency for individual messages
type BatchProducer struct {
	producer sarama.AsyncProducer
	mu       sync.RWMutex
	closed   bool
	wg       sync.WaitGroup
}

type ProducerConfig struct {
	BatchSize   int
	Compression sarama.CompressionCodec
}

func NewBatchProducer(brokers []string, config ProducerConfig) (*BatchProducer, error) {
	saramaConfig := sarama.NewConfig()

	// Optimize for maximum throughput
	saramaConfig.Producer.Return.Successes = false // Change to false for better performance
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Producer.Compression = config.Compression
	saramaConfig.Producer.MaxMessageBytes = 1000000
	saramaConfig.Producer.Flush.MaxMessages = config.BatchSize
	saramaConfig.Producer.Flush.Frequency = 1 * time.Millisecond
	saramaConfig.Producer.RequiredAcks = sarama.NoResponse // Changed for maximum throughput

	// Optimize batch settings
	saramaConfig.Producer.Flush.Bytes = 64 * 1024 // 64KB batch size

	// Channel buffering
	saramaConfig.ChannelBufferSize = 256 * 1024 // Increased buffer size

	// Performance optimizations
	saramaConfig.Producer.Idempotent = false   // Disable idempotence for speed
	saramaConfig.Producer.CompressionLevel = 1 // Fastest compression level
	saramaConfig.Producer.Partitioner = sarama.NewHashPartitioner

	producer, err := sarama.NewAsyncProducer(brokers, saramaConfig)
	if err != nil {
		return nil, err
	}

	bp := &BatchProducer{
		producer: producer,
		closed:   false,
	}

	// Only handle errors since we disabled success returns
	bp.handleAsyncResults()
	return bp, nil
}

func (p *BatchProducer) handleAsyncResults() {
	p.wg.Add(1)
	// Only handle errors
	go func() {
		defer p.wg.Done()
		for err := range p.producer.Errors() {
			if err != nil {
				log.Printf("Failed to send message: %v\n", err.Err)
			}
		}
	}()
}

func (p *BatchProducer) SendMessage(topic, key, value string) error {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return ErrProducerClosed
	}
	p.mu.RUnlock()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	p.producer.Input() <- msg
	return nil
}

func (p *BatchProducer) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	p.mu.Unlock()

	err := p.producer.Close()
	p.wg.Wait()
	return err
}
