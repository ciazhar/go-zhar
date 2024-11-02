package service

import (
	"context"

	"sync"

	"github.com/ciazhar/go-start-small/examples/rabbitmq_publish_consume_testcontainers/internal/model"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/rabbitmq"
)

type BasicService struct {
	ctx      context.Context
	wg       *sync.WaitGroup
	rabbitmq *rabbitmq.RabbitMQ
}

func (e BasicService) StartRabbitConsumer() {
	e.rabbitmq.StartConsumers(e.ctx, []rabbitmq.ConsumerConfig{
		{
			Queue:   model.QueueBasic,
			Handler: e.ConsumeRabbitmq,
		},
	}, e.wg)

}

func (e BasicService) PublishRabbitmq(message string) {
	e.rabbitmq.PublishMessage(context.Background(), model.QueueBasic, message)
}

func (e BasicService) PublishTTLRabbitmq(message string) {
	e.rabbitmq.PublishMessage(context.Background(), model.QueueBasic, message, 1000)
}

func (e BasicService) ConsumeRabbitmq(message string) {
	logger.LogInfo(e.ctx, "Received Basic Message", map[string]interface{}{"message": message})
}

func NewBasicService(
	ctx context.Context,
	rabbitmq *rabbitmq.RabbitMQ,
	wg *sync.WaitGroup,
) *BasicService {
	return &BasicService{
		ctx:      ctx,
		wg:       wg,
		rabbitmq: rabbitmq,
	}
}
