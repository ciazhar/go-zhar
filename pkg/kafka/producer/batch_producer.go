package producer

import (
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type Message struct {
	Topic string
	Key   []byte
	Value []byte
}

type BatchProducer struct {
	producer  sarama.SyncProducer
	input     chan Message
	batch     []*sarama.ProducerMessage
	batchSize int
	done      chan struct{}
	wg        sync.WaitGroup
}

func NewBatchProducer(brokers []string, batchSize int) (*BatchProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &BatchProducer{
		producer:  producer,
		input:     make(chan Message),
		batch:     make([]*sarama.ProducerMessage, 0, batchSize),
		batchSize: batchSize,
		done:      make(chan struct{}),
	}, nil
}

func (p *BatchProducer) Start() {
	log.Println("Starting batch processor")
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for msg := range p.input {
			producerMsg := &sarama.ProducerMessage{
				Topic: msg.Topic,
				Key:   sarama.ByteEncoder(msg.Key),
				Value: sarama.ByteEncoder(msg.Value),
			}

			p.batch = append(p.batch, producerMsg)

			if len(p.batch) >= p.batchSize {
				log.Printf("Batch full, flushing %d messages", len(p.batch))
				p.flush()
			}
		}

		// Flush any remaining messages when the input channel is closed
		log.Println("Messages channel closed, flushing remaining messages")
		if len(p.batch) > 0 {
			p.flush()
		}
		close(p.done)
	}()
}

func (p *BatchProducer) flush() {
	if len(p.batch) == 0 {
		return
	}

	log.Printf("Flushing batch of %d messages", len(p.batch))
	err := p.producer.SendMessages(p.batch)
	if err != nil {
		log.Printf("Failed to send messages: %v", err)
	}

	// Clear the batch
	p.batch = p.batch[:0]
}

func (p *BatchProducer) Input() chan<- Message {
	return p.input
}

func (p *BatchProducer) Close() error {
	close(p.input)
	p.wg.Wait()
	return p.producer.Close()
}

func (p *BatchProducer) Done() <-chan struct{} {
	return p.done
}
