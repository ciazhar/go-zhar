package internal

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/rabbitmq/publish-consume/internal/controller"
	"github.com/ciazhar/go-zhar/examples/rabbitmq/publish-consume/internal/service"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/rabbitmq"
	"github.com/gofiber/fiber/v2"
	"sync"
)

func Init(ctx context.Context, router fiber.Router, mq *rabbitmq.RabbitMQ, wg *sync.WaitGroup, logger logger.Logger) {
	s := service.NewBasicService(ctx, mq, wg, logger)
	c := controller.NewBasicController(s)

	s.StartRabbitConsumer()

	r := router.Group("/basic")
	r.Post("/", c.Publish)
	r.Post("/ttl", c.PublishTTL)
}
