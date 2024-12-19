package producer

import (
	"github.com/IBM/sarama"
	"hash/fnv"
)

// CustomPartitionerProducer implements a producer with custom partitioning logic
// Pros: Fine-grained control over message distribution
// Cons: Need to manage partition assignment logic
type CustomPartitionerProducer struct {
	producer sarama.AsyncProducer
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
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = NewCustomPartitioner

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return &CustomPartitionerProducer{producer: producer}, nil
}
