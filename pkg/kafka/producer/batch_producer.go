package producer

import (
	"github.com/IBM/sarama"
	"sync"
	"time"
)

// BatchProducer BATCH PRODUCER
// Pros: Improved throughput, better network utilization
// Cons: Increased latency for individual messages
type BatchProducer struct {
	Producer     sarama.AsyncProducer
	batchSize    int
	flushTimeout time.Duration
	Messages     chan *sarama.ProducerMessage
	Wg           sync.WaitGroup
}

func NewBatchProducer(brokerList []string, batchSize int, flushTimeout time.Duration) (*BatchProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	// Batch settings
	config.Producer.Flush.Messages = batchSize
	config.Producer.Flush.Frequency = flushTimeout

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	bp := &BatchProducer{
		Producer:     producer,
		batchSize:    batchSize,
		flushTimeout: flushTimeout,
		Messages:     make(chan *sarama.ProducerMessage, batchSize),
	}

	bp.start()
	return bp, nil
}

func (p *BatchProducer) start() {
	p.Wg.Add(1)
	go func() {
		defer p.Wg.Done()
		batch := make([]*sarama.ProducerMessage, 0, p.batchSize)
		timer := time.NewTimer(p.flushTimeout)

		for {
			select {
			case msg, ok := <-p.Messages:
				if !ok {
					// Channel closed, flush remaining messages
					p.flush(batch)
					return
				}

				batch = append(batch, msg)
				if len(batch) >= p.batchSize {
					p.flush(batch)
					batch = make([]*sarama.ProducerMessage, 0, p.batchSize)
					timer.Reset(p.flushTimeout)
				}

			case <-timer.C:
				if len(batch) > 0 {
					p.flush(batch)
					batch = make([]*sarama.ProducerMessage, 0, p.batchSize)
				}
				timer.Reset(p.flushTimeout)
			}
		}
	}()
}

func (p *BatchProducer) flush(batch []*sarama.ProducerMessage) {
	for _, msg := range batch {
		p.Producer.Input() <- msg
	}
}
