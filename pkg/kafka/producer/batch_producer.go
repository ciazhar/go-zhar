package producer

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/IBM/sarama"
)

type Message struct {
	Topic string
	Key   []byte
	Value []byte
}

type BatchProducer struct {
	producer     sarama.SyncProducer
	input        chan Message
	batch        []*sarama.ProducerMessage
	batchSize    int
	done         chan struct{}
	wg           sync.WaitGroup
	messagesSent uint64
	batchesSent  uint64
	errors       uint64
	mu           sync.RWMutex
	closed       bool
}

type ProducerConfig struct {
	BatchSize   int
	Compression sarama.CompressionCodec
}

func NewBatchProducer(brokers []string, config ProducerConfig) (*BatchProducer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Compression = config.Compression

	producer, err := sarama.NewSyncProducer(brokers, saramaConfig)
	if err != nil {
		return nil, err
	}

	// Increase buffer size to handle high throughput
	bufferSize := config.BatchSize * 10 // or another appropriate multiplier

	return &BatchProducer{
		producer:  producer,
		input:     make(chan Message, bufferSize),
		batch:     make([]*sarama.ProducerMessage, 0, config.BatchSize),
		batchSize: config.BatchSize,
		done:      make(chan struct{}),
	}, nil
}

func (p *BatchProducer) SendMessage(topic, key, value string) error {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return ErrProducerClosed
	}
	p.mu.RUnlock()

	// Blocking send
	p.input <- Message{
		Topic: topic,
		Key:   []byte(key),
		Value: []byte(value),
	}
	return nil
}

// Make the processing loop more efficient
func (p *BatchProducer) Start() {
	//log.Println("Starting batch processor")
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		ticker := time.NewTicker(100 * time.Millisecond) // Flush interval
		defer ticker.Stop()

		for {
			select {
			case msg, ok := <-p.input:
				if !ok {
					// Channel closed, flush remaining messages and exit
					if len(p.batch) > 0 {
						p.flush()
					}
					close(p.done)
					return
				}

				producerMsg := &sarama.ProducerMessage{
					Topic: msg.Topic,
					Key:   sarama.StringEncoder(msg.Key),
					Value: sarama.StringEncoder(msg.Value),
				}

				p.batch = append(p.batch, producerMsg)

				if len(p.batch) >= p.batchSize {
					p.flush()
				}

			case <-ticker.C:
				// Periodically flush if there are messages
				if len(p.batch) > 0 {
					p.flush()
				}
			}
		}
	}()
}

func (p *BatchProducer) flush() {
	if len(p.batch) == 0 {
		return
	}

	err := p.producer.SendMessages(p.batch)
	if err != nil {
		atomic.AddUint64(&p.errors, 1)
		log.Printf("Failed to send messages: %v", err)
	} else {
		atomic.AddUint64(&p.messagesSent, uint64(len(p.batch)))
		atomic.AddUint64(&p.batchesSent, 1)
	}

	p.batch = p.batch[:0]
}

func (p *BatchProducer) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	p.mu.Unlock()

	close(p.input)
	p.wg.Wait()
	return p.producer.Close()
}
