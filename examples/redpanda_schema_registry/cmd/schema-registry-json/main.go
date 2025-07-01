package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/IBM/sarama"
	"github.com/linkedin/goavro"
)

// Schema Registry URL
const schemaRegistryURL = "http://localhost:38081"

// Avro Schema for User
const avroSchema = `{
    "type": "record",
    "name": "User",
    "fields": [
        {"name": "id", "type": "int"},
        {"name": "name", "type": "string"},
        {"name": "email", "type": "string"}
    ]
}`

// Fetch Schema ID from Schema Registry
func getSchemaID(schema string) int {
	reqBody := fmt.Sprintf(`{"schema": "%s"}`, strings.ReplaceAll(schema, `"`, `\"`))
	resp, err := http.Post(schemaRegistryURL+"/subjects/users-value/versions", "application/json", strings.NewReader(reqBody))
	if err != nil {
		log.Fatal("Schema registration failed:", err)
	}
	defer resp.Body.Close()

	var result struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal("Failed to parse Schema ID:", err)
	}
	return result.ID
}

// Serialize Avro with Schema ID
func serializeAvro(schemaID int, user map[string]interface{}) ([]byte, error) {
	// Create an Avro codec
	codec, err := goavro.NewCodec(avroSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to create Avro codec: %v", err)
	}

	// Convert native Go map to Avro binary format
	binaryData, err := codec.BinaryFromNative(nil, user)
	if err != nil {
		return nil, fmt.Errorf("failed to encode Avro: %v", err)
	}

	// Prepend the schema ID (4 bytes)
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.BigEndian, int32(schemaID)) // Write Schema ID
	buf.Write(binaryData)                                     // Append encoded Avro data

	return buf.Bytes(), nil
}

// Produce Message to Kafka
func produceMessage() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	brokers := []string{"localhost:19092", "localhost:29092", "localhost:39092"}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatal("Failed to start producer:", err)
	}
	defer producer.Close()

	// Fetch Schema ID from Schema Registry
	schemaID := 3

	log.Println("Schema ID:", schemaID)

	// User Data
	user := map[string]interface{}{
		"id":    1,
		"name":  "John Doe",
		"email": "john@example.com",
	}

	// Serialize with Avro
	avroData, err := serializeAvro(schemaID, user)
	if err != nil {
		log.Fatal("Failed to serialize Avro data:", err)
	}

	// Send to Redpanda
	msg := &sarama.ProducerMessage{
		Topic: "users-avro-6",
		Value: sarama.ByteEncoder(avroData),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		log.Fatal("Failed to send message:", err)
	}

	fmt.Println("Produced Avro message with Schema Registry!")
}

func main() {
	produceMessage()
}
