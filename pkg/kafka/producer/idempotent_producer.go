package producer

import "github.com/IBM/sarama"

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
