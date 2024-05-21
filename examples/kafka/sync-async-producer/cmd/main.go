package main

import (
	kafka2 "github.com/ciazhar/go-zhar/examples/kafka/sync-async-producer/internal/kafka"
	"github.com/ciazhar/go-zhar/examples/kafka/sync-async-producer/internal/kafka/model"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/kafka"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {

	// Logger
	log := logger.Init(logger.Config{
		ConsoleLoggingEnabled: true,
	})

	// Environment configuration
	env.Init("config.json", log)

	// Kafka configuration
	admin := kafka.NewAdmin(viper.GetString("kafka.brokers"), log)
	admin.CreateTopic(model.TopicSync, kafka.CreateTopicConfig{
		NumPartitions:     3,
		ReplicationFactor: 1,
	})
	admin.CreateTopic(model.TopicAsync, kafka.CreateTopicConfig{
		NumPartitions:     3,
		ReplicationFactor: 1,
	})

	syncProducer := kafka.NewSyncProducer(viper.GetString("kafka.brokers"), log)
	asyncProducer := kafka.NewAsyncProducer(viper.GetString("kafka.brokers"), log)

	app := fiber.New()
	kafka2.Init(app, syncProducer, asyncProducer)

	err := app.Listen(":" + viper.GetString("application.port"))
	if err != nil {
		log.Fatalf("fiber failed to start : %v", err)
	}
}
