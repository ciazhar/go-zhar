package producer

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

// SyncProducer SYNCHRONOUS PRODUCER
// Pros: Simple to use, immediate feedback on success/failure
// Cons: Lower throughput due to waiting for response
type SyncProducer struct {
	producer sarama.SyncProducer
}

func NewSyncProducer(brokerList []string, maxRetry int) (*SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = maxRetry

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return &SyncProducer{producer: producer}, nil
}

func (p *SyncProducer) SendMessage(topic, key string, value any) (partition int32, offset int64, err error) {

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
