package main

import (
	"fmt"
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

	consumer := kafka.NewConsumer(viper.GetString("kafka.brokers"), log)
	defer consumer.Close()

	admin := kafka.NewAdmin(viper.GetString("kafka.brokers"), log)
	defer admin.Close()

	admin.CreateTopic(model.TopicName, kafka.CreateTopicConfig{
		NumPartitions:     3,
		ReplicationFactor: 1,
	})
	admin.CreateTopic(model.Topic2Name, kafka.CreateTopicConfig{
		NumPartitions:     3,
		ReplicationFactor: 1,
	})

	KafkaConsumer(consumer)

	forever := make(chan bool)
	<-forever
}

func KafkaConsumer(consumer kafka.Consumer) {
	go consumer.ConsumeMessages(model.TopicName, HandleMessage)
	go consumer.ConsumeMessages(model.Topic2Name, HandleMessage2)
}

func HandleMessage(string2 string) {
	fmt.Printf("From Topic 1. Received a message: %s\n", string2)
}

func HandleMessage2(string2 string) {
	fmt.Printf("From Topic 2. Received a message: %s\n", string2)
}
