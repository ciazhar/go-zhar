package producer

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

// IdempotentProducer IDEMPOTENT PRODUCER
// Pros: Exactly-once delivery semantics
// Cons: Higher latency, requires newer Kafka version
type IdempotentProducer struct {
	producer sarama.SyncProducer
}

func NewIdempotentProducer(brokerList []string) (*IdempotentProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Idempotent = true
	config.Producer.Return.Successes = true
	config.Net.MaxOpenRequests = 1 // Required for idempotence

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return &IdempotentProducer{producer: producer}, nil
}

func (p *IdempotentProducer) SendMessage(topic, key string, value any) (partition int32, offset int64, err error) {

	marshal, err := json.Marshal(value)
	if err != nil {
		return 0, 0, err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(marshal), // use string encoder for text and json data, for binary data use byte encoder
	}

	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}

	return p.producer.SendMessage(msg)
}
