package service

import (
	"context"

	"sync"

	"github.com/ciazhar/go-start-small/examples/rabbitmq_publish_consume_testcontainers/internal/model"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/rabbitmq"
)

type NotificationService struct {
	ctx      context.Context
	wg       *sync.WaitGroup
	rabbitmq *rabbitmq.RabbitMQ
}

func (e NotificationService) StartRabbitConsumer() {
	e.rabbitmq.StartConsumers(e.ctx, []rabbitmq.ConsumerConfig{
		{
			Queue:   model.OrderStatusQueue,
			Handler: e.ConsumeOrderStatusQueue,
		},
		{
			Queue:   model.PaymentReminderQueue,
			Handler: e.ConsumePaymentReminderQueue,
		},
	}, e.wg)

}

func (e NotificationService) PublishRabbitmq(message string) {
	e.rabbitmq.PublishMessage(context.Background(), model.OrderStatusQueue, message)
}

func (e NotificationService) PublishTTLRabbitmq(message string) {
	e.rabbitmq.PublishMessage(context.Background(), model.PaymentReminderQueue, message, 1000)
}

func (e NotificationService) ConsumeOrderStatusQueue(message string) {
	logger.LogInfo(e.ctx, "Received Basic Message", map[string]interface{}{"message": message})
	e.SendNotificationViaSMSAndEmail(message)
}

func (e NotificationService) ConsumePaymentReminderQueue(message string) {
	logger.LogInfo(e.ctx, "Received Basic Message", map[string]interface{}{"message": message})
}

func (e NotificationService) SendNotificationViaSMSAndEmail(message string) {
	logger.LogInfo(e.ctx, "Sending notification via SMS and Email", map[string]interface{}{"message": message})
}

func NewBasicService(
	ctx context.Context,
	rabbitmq *rabbitmq.RabbitMQ,
	wg *sync.WaitGroup,
) *NotificationService {
	return &NotificationService{
		ctx:      ctx,
		wg:       wg,
		rabbitmq: rabbitmq,
	}
}
