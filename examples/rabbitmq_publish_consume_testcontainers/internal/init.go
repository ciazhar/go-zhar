package internal

import (
	"context"

	"sync"

	"github.com/ciazhar/go-start-small/examples/rabbitmq_publish_consume_testcontainers/internal/controller"
	"github.com/ciazhar/go-start-small/examples/rabbitmq_publish_consume_testcontainers/internal/service"
	"github.com/ciazhar/go-start-small/pkg/rabbitmq"
	"github.com/gofiber/fiber/v2"
)

func Init(ctx context.Context, router fiber.Router, mq *rabbitmq.RabbitMQ, wg *sync.WaitGroup) {
	s := service.NewBasicService(ctx, mq, wg)
	c := controller.NewBasicController(s)

	s.StartRabbitConsumer()

	r := router.Group("/notification")
	r.Post("/order-status", c.UpdateOrderStatus)
	r.Post("/payment-reminder", c.SendPaymentReminder)
}
