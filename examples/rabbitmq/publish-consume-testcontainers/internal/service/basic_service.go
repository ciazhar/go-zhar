package service

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/rabbitmq/publish-consume-testcontainers/internal/model"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/rabbitmq"
	"sync"
)

type BasicService struct {
	ctx      context.Context
	wg       *sync.WaitGroup
	logger   *logger.Logger
	rabbitmq *rabbitmq.RabbitMQ
}

func (e BasicService) StartRabbitConsumer() {
	e.rabbitmq.StartConsumers(e.ctx, []rabbitmq.ConsumerConfig{
		{
			Queue:   model.QueueBasic,
			Handler: e.ConsumeRabbitmq,
		},
	}, e.wg, e.logger)

}

func (e BasicService) PublishRabbitmq(message string) {
	e.rabbitmq.PublishMessage(context.Background(), model.QueueBasic, message)
}

func (e BasicService) PublishTTLRabbitmq(message string) {
	e.rabbitmq.PublishMessageWithTTL(context.Background(), model.QueueBasic, message, 1000)
}

func (e BasicService) ConsumeRabbitmq(message string) {
	e.logger.Infof("Received Basic Message: %s", message)
}

func NewBasicService(
	ctx context.Context,
	rabbitmq *rabbitmq.RabbitMQ,
	wg *sync.WaitGroup,
	logger *logger.Logger,
) *BasicService {
	return &BasicService{
		ctx:      ctx,
		wg:       wg,
		logger:   logger,
		rabbitmq: rabbitmq,
	}
}
