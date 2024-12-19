package producer

import "github.com/IBM/sarama"

// SyncProducer SYNCHRONOUS PRODUCER
// Pros: Simple to use, immediate feedback on success/failure
// Cons: Lower throughput due to waiting for response
type SyncProducer struct {
	producer sarama.SyncProducer
}

func NewSyncProducer(brokerList []string) (*SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return &SyncProducer{producer: producer}, nil
}

func (p *SyncProducer) SendMessage(topic, key, value string) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	return p.producer.SendMessage(msg)
}
