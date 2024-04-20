package event

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/kafka/custom-consumer-group/internal/event/service"
	"github.com/ciazhar/go-zhar/pkg/kafka"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/spf13/viper"
)

func Init(log logger.Logger, producer *kafka.SyncProducer) {

	s := service.NewEventService(log, producer)

	kafka.StartConsumers(context.Background(), viper.GetString("kafka.brokers"), map[string]kafka.ConsumerConfig{
		"event-group": {
			Topics:  []string{"event"},
			Handler: s,
		},
	}, log)
}
