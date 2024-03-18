package basic

import (
	"github.com/ciazhar/go-zhar/examples/message-broker/rabbitmq/basic/internal/basic/controller"
	"github.com/ciazhar/go-zhar/examples/message-broker/rabbitmq/basic/internal/basic/model"
	"github.com/ciazhar/go-zhar/examples/message-broker/rabbitmq/basic/internal/basic/service"
	"github.com/ciazhar/go-zhar/pkg/message_broker/rabbitmq"
	"github.com/gofiber/fiber/v2"
)

func Init(router fiber.Router, mq *rabbitmq.RabbitMQ, rabbitmqChan chan struct{}) {
	s := service.NewBasicService(mq)
	c := controller.NewBasicController(s)

	go mq.ConsumeMessages(model.QueueBasic, s.ConsumeRabbitmq, rabbitmqChan)

	r := router.Group("/basic")
	r.Post("/", c.Publish)
}
