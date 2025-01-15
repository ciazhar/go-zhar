package producer

import (
	"log"
	"sync"
	"sync/atomic"

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

	return &BatchProducer{
		producer:  producer,
		input:     make(chan Message),
		batch:     make([]*sarama.ProducerMessage, 0, config.BatchSize),
		batchSize: config.BatchSize,
		done:      make(chan struct{}),
	}, nil
}

func (p *BatchProducer) Start() {
	//log.Println("Starting batch processor")
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
				p.flush()
			}
		}

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

func (p *BatchProducer) Stats() (uint64, uint64, uint64) {
	return atomic.LoadUint64(&p.messagesSent),
		atomic.LoadUint64(&p.batchesSent),
		atomic.LoadUint64(&p.errors)
}
