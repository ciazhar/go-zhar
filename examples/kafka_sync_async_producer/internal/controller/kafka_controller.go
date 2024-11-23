package Controller

import (
	"github.com/IBM/sarama"
	"github.com/ciazhar/go-start-small/examples/kafka_sync_async_producer/internal/model"
	"github.com/ciazhar/go-start-small/pkg/kafka"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	syncConsumer  sarama.SyncProducer
	asyncConsumer sarama.AsyncProducer
}

func NewController(
	syncConsumer sarama.SyncProducer,
	asyncConsumer sarama.AsyncProducer,
) *Controller {
	return &Controller{
		syncConsumer:  syncConsumer,
		asyncConsumer: asyncConsumer,
	}
}

func (c *Controller) SyncProducer(ctx *fiber.Ctx) error {
	text := "Message Sent To Sync Topic!"
	err := kafka.SendMessage(c.syncConsumer, model.TopicSync, text)
	if err != nil {
		return err
	}
	return ctx.SendString(text)
}

func (c *Controller) AsyncProducer(ctx *fiber.Ctx) error {
	text := "Message Sent To Async Topic!"
	kafka.SendAsyncMessage(c.asyncConsumer, model.TopicAsync, text)
	return ctx.SendString(text)
}
