package main

import (
	"github.com/ciazhar/go-zhar/examples/kafka/simple-consumer-producer/model"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/kafka"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/spf13/viper"
)

func main() {

	// Logger
	log := logger.Init()

	// Environment configuration
	env.Init("config.json", log)

	producer := kafka.NewSyncProducer(viper.GetString("kafka.brokers"), log)
	defer producer.Close()

	producer.PublishMessage(model.TopicName, "Hello, Kafka!")
	producer.PublishMessage(model.Topic2Name, "Hello, Kafka!")
}
