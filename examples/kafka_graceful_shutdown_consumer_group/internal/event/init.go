package event

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/kafka_graceful_shutdown_consumer_group/internal/event/service"
	"github.com/ciazhar/go-start-small/pkg/kafka/consumer"
	"github.com/spf13/viper"
	"sync"
)

func Init(ctx context.Context, wg *sync.WaitGroup) {

	s := service.NewEventService()

	consumer.StartConsumerGroup(
		ctx,
		viper.GetString("kafka.brokers"),
		map[string]consumer.ConsumerGroup{
			"my-topic": {
				GroupID: "my-group",
				Topic:   "my-topic",
				Process: s.ProcessEvent,
			},
		},
		wg,
		viper.GetString("kafka.assignor"),
		viper.GetBool("kafka.offsetOldest"),
	)
}
