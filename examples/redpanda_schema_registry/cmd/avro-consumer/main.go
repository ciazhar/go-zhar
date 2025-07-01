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

// Avro Schema
const avroSchema = `{
    "type": "record",
    "name": "user",
    "fields": [
        {"name": "id", "type": "int"},
        {"name": "name", "type": "string"},
        {"name": "email", "type": "string"}
    ]
}`

// Serialize Avro with Schema ID
func serializeAvro(schemaID int, user map[string]interface{}) []byte {
	codec, _ := goavro.NewCodec(avroSchema)
	binaryData, _ := codec.BinaryFromNative(nil, user)

	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.BigEndian, int32(schemaID))
	buf.Write(binaryData)

	return buf.Bytes()
}

// Produce Message to Kafka
func produceMessage() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	brokers := []string{"localhost:19092", "localhost:29092", "localhost:39092"}
	producer, _ := sarama.NewSyncProducer(brokers, config)
	defer producer.Close()

	// Fetch Schema ID from Schema Registry
	schemaID := 2

	user := map[string]interface{}{
		"id":    1,
		"name":  "John Doe",
		"email": "john@example.com",
	}

	avroData := serializeAvro(schemaID, user)
	msg := &sarama.ProducerMessage{
		Topic: "benchmark-avro-value-2",
		Value: sarama.ByteEncoder(avroData),
	}

	_, _, _ = producer.SendMessage(msg)
	fmt.Println("Produced Avro message with Schema Registry!")
}

func main() {
	produceMessage()
}

//
