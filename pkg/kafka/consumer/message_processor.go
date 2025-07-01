package consumer

import "github.com/IBM/sarama"

// MessageProcessor is a function type for processing messages.
type MessageProcessor func(msg *sarama.ConsumerMessage) error
