package basic

import (
	"github.com/ciazhar/go-zhar/examples/rabbitmq/clean-architecture/internal/basic/controller"
	"github.com/ciazhar/go-zhar/examples/rabbitmq/clean-architecture/internal/basic/model"
	"github.com/ciazhar/go-zhar/examples/rabbitmq/clean-architecture/internal/basic/service"
	"github.com/ciazhar/go-zhar/pkg/rabbitmq"
	"github.com/gofiber/fiber/v2"
)

func Init(router fiber.Router, mq *rabbitmq.RabbitMQ, rabbitmqChan chan struct{}) {
	s := service.NewBasicService(mq)
	c := controller.NewBasicController(s)

	go mq.ConsumeMessages(model.QueueBasic, s.ConsumeRabbitmq, rabbitmqChan)

	r := router.Group("/basic")
	r.Post("/", c.Publish)
}
