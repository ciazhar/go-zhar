package service

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/rabbitmq/publish-consume/internal/model"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/rabbitmq"
	"sync"
)

type BasicService interface {
	PublishRabbitmq(message string)
	PublishTTLRabbitmq(message string)
	ConsumeRabbitmq(message string)
	StartRabbitConsumer()
}

type basicService struct {
	ctx      context.Context
	wg       *sync.WaitGroup
	logger   logger.Logger
	rabbitmq *rabbitmq.RabbitMQ
}

func (e basicService) StartRabbitConsumer() {
	e.rabbitmq.StartConsumers(e.ctx, []rabbitmq.ConsumerConfig{
		{
			Queue:   model.QueueBasic,
			Handler: e.ConsumeRabbitmq,
		},
	}, e.wg, e.logger)

}

func (e basicService) PublishRabbitmq(message string) {
	e.rabbitmq.PublishMessage(context.Background(), model.QueueBasic, message)
}

func (e basicService) PublishTTLRabbitmq(message string) {
	e.rabbitmq.PublishMessageWithTTL(context.Background(), model.QueueBasic, message, 1000)
}

func (e basicService) ConsumeRabbitmq(message string) {
	e.logger.Infof("Received Basic Message: %s", message)
}

func NewBasicService(
	ctx context.Context,
	rabbitmq *rabbitmq.RabbitMQ,
	wg *sync.WaitGroup,
	logger logger.Logger,
) BasicService {
	return basicService{
		ctx:      ctx,
		wg:       wg,
		logger:   logger,
		rabbitmq: rabbitmq,
	}
}
