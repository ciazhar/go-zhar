package event

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/kafka/consumer-group/internal/event/service"
	"github.com/ciazhar/go-zhar/pkg/kafka"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/spf13/viper"
	"sync"
)

func Init(log *logger.Logger, ctx context.Context, wg *sync.WaitGroup) {

	s := service.NewEventService(log)

	kafka.StartConsumers(ctx, viper.GetString("kafka.brokers"), map[string]kafka.ConsumerConfig{
		"event-group": {
			Topics:  []string{"event"},
			Handler: kafka.NewBasicConsumerHandler(s.ProcessEvent),
		},
	}, wg, log)
}
