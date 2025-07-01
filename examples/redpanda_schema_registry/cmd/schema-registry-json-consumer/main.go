package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/riferrei/srclient"
)

// User struct
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Schema Registry URL
const schemaRegistryURL = "http://localhost:8081"

func consumeWithSchemaRegistry(topic string) {
	// Kafka Consumer Configuration
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Connect to Kafka
	brokers := []string{"localhost:9092"}
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	// Connect to Schema Registry
	client := srclient.CreateSchemaRegistryClient(schemaRegistryURL)

	// Consume Messages
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal(err)
	}
	defer partitionConsumer.Close()

	fmt.Println("Listening for messages...")

	for msg := range partitionConsumer.Messages() {
		// Extract Schema ID
		if len(msg.Value) < 4 {
			log.Println("Invalid message: too short")
			continue
		}
		schemaID := int(binary.BigEndian.Uint32(msg.Value[:4]))

		// Get Schema from Registry
		schema, err := client.GetSchema(schemaID)
		if err != nil {
			log.Fatal("Schema lookup failed:", err)
		}

		// Deserialize JSON Payload
		jsonData := msg.Value[4:] // Skip schema ID prefix
		var user User
		if err := json.Unmarshal(jsonData, &user); err != nil {
			log.Println("Failed to unmarshal JSON:", err)
			continue
		}

		// Print Message
		fmt.Printf("Consumed: %+v\n", user)
	}
}

func main() {
	consumeWithSchemaRegistry("users-json")
}
