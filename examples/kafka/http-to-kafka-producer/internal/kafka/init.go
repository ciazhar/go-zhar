package kafka

import (
	controller2 "github.com/ciazhar/go-zhar/examples/kafka/http-to-kafka-producer/internal/kafka/controller"
	"github.com/ciazhar/go-zhar/pkg/kafka"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App,
	syncConsumer *kafka.SyncProducer,
	asyncConsumer *kafka.AsyncProducer) {
	controller := controller2.NewController(syncConsumer, asyncConsumer)

	r := app.Group("/")
	r.Get("/async", controller.AsyncProducer)
	r.Get("/sync", controller.SyncProducer)
}
