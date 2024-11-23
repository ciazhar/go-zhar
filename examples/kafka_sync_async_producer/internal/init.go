package internal

import (
	"github.com/IBM/sarama"
	Controller "github.com/ciazhar/go-start-small/examples/kafka_sync_async_producer/internal/controller"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App,
	syncConsumer sarama.SyncProducer,
	asyncConsumer sarama.AsyncProducer) {
	controller := Controller.NewController(syncConsumer, asyncConsumer)

	r := app.Group("/")
	r.Get("/async", controller.AsyncProducer)
	r.Get("/sync", controller.SyncProducer)
}
