package producer

import (
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
}

func NewAsyncProducer(brokerList []string) (*AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	ap := &AsyncProducer{producer: producer}
	ap.handleAsyncResults()
	return ap, nil
}

func (p *AsyncProducer) handleAsyncResults() {
	p.wg.Add(2)
	// Success handler
	go func() {
		defer p.wg.Done()
		for success := range p.producer.Successes() {
			log.Printf("Message sent successfully: topic=%s partition=%d offset=%d\n",
				success.Topic, success.Partition, success.Offset)
		}
	}()

	// Error handler
	go func() {
		defer p.wg.Done()
		for err := range p.producer.Errors() {
			log.Printf("Failed to send message: %v\n", err)
		}
	}()
}

func (p *AsyncProducer) SendMessage(topic, key, value string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}
	p.producer.Input() <- msg
}
