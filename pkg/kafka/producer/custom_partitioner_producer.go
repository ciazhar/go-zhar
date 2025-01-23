package producer

import (
	"github.com/IBM/sarama"
	"hash/fnv"
	"log"
	"sync"
	"time"
)

// CustomPartitionerProducer implements a producer with custom partitioning logic
// Pros: Fine-grained control over message distribution
// Cons: Need to manage partition assignment logic
type CustomPartitionerProducer struct {
	producer sarama.AsyncProducer
	mu       sync.RWMutex
	closed   bool
	wg       sync.WaitGroup
}

type CustomPartitioner struct {
	partition int32
}

// NewCustomPartitioner creates a new instance of CustomPartitioner
func NewCustomPartitioner(topic string) sarama.Partitioner {
	return &CustomPartitioner{}
}

func (p *CustomPartitioner) Partition(message *sarama.ProducerMessage, numPartitions int32) (int32, error) {
	// Example: Partition based on message key hash
	if message.Key == nil {
		return 0, nil
	}

	// Convert key to bytes for hashing
	var keyBytes []byte
	var err error

	switch k := message.Key.(type) {
	case sarama.StringEncoder:
		keyBytes, err = k.Encode()
		if err != nil {
			return 0, err
		}
	case sarama.ByteEncoder:
		keyBytes = k
	default:
		// If key is neither string nor byte encoder, use partition 0
		return 0, nil
	}

	hash := fnv.New32a()
	hash.Write(keyBytes)
	return int32(hash.Sum32()) % numPartitions, nil
}

func (p *CustomPartitioner) RequiresConsistency() bool {
	return true
}

func NewCustomPartitionerProducer(brokerList []string) (*CustomPartitionerProducer, error) {
	config := sarama.NewConfig()

	// Optimize for throughput
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = true
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.MaxMessages = 1000
	config.Producer.Flush.Frequency = 1 * time.Millisecond
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Flush.Bytes = 64 * 1024
	config.ChannelBufferSize = 256 * 1024
	config.Producer.MaxMessageBytes = 1000000
	config.Producer.CompressionLevel = 1

	// Set custom partitioner
	config.Producer.Partitioner = NewCustomPartitioner

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	cpp := &CustomPartitionerProducer{
		producer: producer,
		closed:   false,
	}

	cpp.handleAsyncResults()
	return cpp, nil
}

func (p *CustomPartitionerProducer) handleAsyncResults() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for err := range p.producer.Errors() {
			if err != nil {
				log.Printf("Failed to send message: %v\n", err.Err)
			}
		}
	}()
}

func (p *CustomPartitionerProducer) SendMessage(topic, key, value string) error {
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

func (p *CustomPartitionerProducer) Close() error {
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
