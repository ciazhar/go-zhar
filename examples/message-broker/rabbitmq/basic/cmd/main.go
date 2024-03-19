package main

import (
	"fmt"
	"github.com/ciazhar/go-zhar/examples/message-broker/rabbitmq/basic/internal/basic"
	"github.com/ciazhar/go-zhar/examples/message-broker/rabbitmq/basic/internal/basic/model"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/message_broker/rabbitmq"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// Channel to receive OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Environment configuration
	env.Init("config.json")

	// RabbitMQ configuration
	ra := rabbitmq.New(viper.GetString("rabbitmq.username"), viper.GetString("rabbitmq.password"), viper.GetString("rabbitmq.host"), viper.GetString("rabbitmq.port"))
	ra.CreateQueue(model.QueueBasic)
	rabbitmqConsumerDone := make(chan struct{})
	defer ra.Close()

	// Fiber configuration
	app := fiber.New()

	// Module initialization
	basic.Init(app, ra, rabbitmqConsumerDone)

	// Start Fiber
	go func() {
		err := app.Listen(":" + viper.GetString("application.port"))
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for termination signal
	<-sigCh
	fmt.Println("Received termination signal. Initiating shutdown...")

	// Initiate shutdown
	close(rabbitmqConsumerDone)
	err := app.Shutdown()
	if err != nil {
		log.Fatal(err)
	}

	// Wait for consumers to finish
	for i := 0; i < 1; i++ { // Adjust the number based on your actual consumers
		select {
		case <-rabbitmqConsumerDone:
			fmt.Printf("Consumer %d has gracefully stopped.\n", i)
		case <-time.After(5 * time.Second):
			fmt.Printf("Consumer %d did not stop in time. Forcefully exiting.\n", i)
		}
	}
}
