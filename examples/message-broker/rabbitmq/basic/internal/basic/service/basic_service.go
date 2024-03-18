package service

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/message-broker/rabbitmq/basic/model"
	"github.com/ciazhar/go-zhar/pkg/message_broker/rabbitmq"
	"github.com/gofiber/fiber/v2/log"
)

type BasicService interface {
	PublishRabbitmq(message string)
	PublishTTLRabbitmq(message string)
	ConsumeRabbitmq(message string)
}

type basicService struct {
	rabbitmq *rabbitmq.RabbitMQ
}

func (e basicService) PublishRabbitmq(message string) {
	e.rabbitmq.PublishMessage(context.Background(), model.QueueBasic, message)
}

func (e basicService) PublishTTLRabbitmq(message string) {
	e.rabbitmq.PublishMessageWithTTL(context.Background(), model.QueueBasic, message, 1000)
}

func (e basicService) ConsumeRabbitmq(message string) {
	log.Info("Received Basic Message: ", message)
}

func NewBasicService(rabbitmq *rabbitmq.RabbitMQ) BasicService {
	return basicService{
		rabbitmq: rabbitmq,
	}
}
