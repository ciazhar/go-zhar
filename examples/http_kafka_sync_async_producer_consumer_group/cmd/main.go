package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/ciazhar/go-start-small/pkg/kafka"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	syncProducer  sarama.SyncProducer
	asyncProducer sarama.AsyncProducer
}

func newController(syncProducer sarama.SyncProducer, asyncProducer sarama.AsyncProducer) *Controller {
	return &Controller{
		syncProducer:  syncProducer,
		asyncProducer: asyncProducer,
	}
}

func processMessage(msg *sarama.ConsumerMessage) error {
	logger.LogInfo(context.Background(), "Processing message", map[string]interface{}{
		"key":   string(msg.Key),
		"value": string(msg.Value),
	})
	return nil // Add your message processing logic here
}

func (a *Controller) sendSyncMessage(c *fiber.Ctx) error {
	message := c.Query("message")
	if message == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Message is required"})
	}

	kafka.SendMessage(a.syncProducer, "my-topic", message)
	return c.JSON(fiber.Map{"status": "success", "message": message})
}

func (a *Controller) sendAsyncMessage(c *fiber.Ctx) error {
	message := c.Query("message")
	if message == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Message is required"})
	}

	kafka.SendAsyncMessage(a.asyncProducer, "my-topic", message)
	return c.JSON(fiber.Map{"status": "success", "message": message})
}

func startConsumer(ctx context.Context, wg *sync.WaitGroup, kafkaConsumer *kafka.KafkaConsumer) {
	defer wg.Done()

	wg.Add(1)
	kafkaConsumer.ConsumeMessages(ctx, wg)
}

func main() {
	// Initialize Fiber app
	app := fiber.New()

	// Kafka configuration
	brokers := []string{"localhost:9092"}

	// Create waitgroup
	var wg sync.WaitGroup

	// Create producers
	syncProducer := kafka.CreateProducer(brokers)
	asyncProducer := kafka.CreateAsyncProducer(brokers)

	// Create consumer
	kafkaConsumer, err := kafka.NewKafkaConsumer(brokers, "my-group", "my-topic", processMessage)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to create Kafka consumer", nil)
	}

	// Create controller
	controller := newController(syncProducer, asyncProducer)

	// Define routes
	app.Post("/send/sync", controller.sendSyncMessage)
	app.Post("/send/async", controller.sendAsyncMessage)

	// Start Kafka consumer in a separate goroutine
	ctx, cancel := context.WithCancel(context.Background())
	go startConsumer(ctx, &wg, kafkaConsumer)

	// Graceful shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		cancel() // Cancel the context to stop the consumer
		app.Shutdown()
	}()

	// Start server
	err = app.Listen(":3000")
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to start server", nil)
	}

	// Wait for the consumer to finish processing messages
	wg.Wait()

	// Cleanup
	syncProducer.Close()
	asyncProducer.Close()
	kafkaConsumer.Close()
}
