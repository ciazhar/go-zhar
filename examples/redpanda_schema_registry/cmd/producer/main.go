package main

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/ciazhar/go-start-small/examples/redpanda_schema_registry/internal/model"
	"github.com/ciazhar/go-start-small/examples/redpanda_schema_registry/proto"
	"github.com/linkedin/goavro"
	proto2 "google.golang.org/protobuf/proto"
	"log"
)

// Serialize JSON
func serializeJSON(user model.User) []byte {
	data, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("Error serializing JSON: %v", err)
	}
	return data
}

// Serialize Avro
func serializeAvro(user model.User) []byte {
	codec, err := goavro.NewCodec(model.AvroSchema)
	if err != nil {
		log.Fatal(err)
	}
	binary, err := codec.BinaryFromNative(nil, map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
	if err != nil {
		log.Fatal(err)
	}
	return binary
}

// Serialize Protobuf
func serializeProtobuf(user model.User) []byte {
	userProto := &proto.UserProto{
		Id:    int32(user.ID),
		Name:  user.Name,
		Email: user.Email,
	}
	data, err := proto2.Marshal(userProto)
	if err != nil {
		log.Fatalf("Error serializing Protobuf: %v", err)
	}
	return data
}

// Produce message to Kafka
func produceMessage(topic string, value []byte) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	brokers := []string{"localhost:19092", "localhost:29092", "localhost:39092"}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
}

func main() {
	user := model.User{ID: 1, Name: "John Doe", Email: "john@example.com"}

	jsonData := serializeJSON(user)
	avroData := serializeAvro(user)
	protoData := serializeProtobuf(user)

	fmt.Printf("JSON Size: %d bytes\n", len(jsonData))
	fmt.Printf("Avro Size: %d bytes\n", len(avroData))
	fmt.Printf("Protobuf Size: %d bytes\n", len(protoData))

	// Send to Kafka
	produceMessage("benchmark-json", jsonData)
	produceMessage("benchmark-avro", avroData)
	produceMessage("benchmark-protobuf", protoData)
}
