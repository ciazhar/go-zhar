package producer

import (
	"errors"
	"github.com/IBM/sarama"
	"log"
	"sync"
	"time"
)

// AsyncProducer ASYNC PRODUCER WITH COMPRESSION
// Pros: Higher throughput, better network utilization
// Cons: More complex error handling, may lose messages if not handled properly
type AsyncProducer struct {
	producer sarama.AsyncProducer
	wg       sync.WaitGroup
	closed   bool
	mu       sync.Mutex
}

func NewAsyncProducer(brokerList []string) (*AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	// Set additional safety configurations
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all replicas
	config.Producer.Retry.Max = 5                    // Retry up to 5 times
	config.Producer.Retry.Backoff = 100 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	ap := &AsyncProducer{
		producer: producer,
		closed:   false,
	}
	ap.handleAsyncResults()
	return ap, nil
}

func (p *AsyncProducer) handleAsyncResults() {
	p.wg.Add(2)
	// Success handler
	go func() {
		defer p.wg.Done()
		for range p.producer.Successes() {
		}
	}()

	// Error handler
	go func() {
		defer p.wg.Done()
		for err := range p.producer.Errors() {
			if err != nil {
				log.Printf("Failed to send message: %v\n", err.Err)
			}
		}
	}()
}

func (p *AsyncProducer) SendMessage(topic, key, value string) error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return ErrProducerClosed
	}
	p.mu.Unlock()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}
	p.producer.Input() <- msg
	return nil
}

// Close gracefully shuts down the producer
func (p *AsyncProducer) Close() error {
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

// ErrProducerClosed is returned when trying to send a message to a closed producer
var ErrProducerClosed = errors.New("producer is closed")
