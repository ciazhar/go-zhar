package main

import (
	"encoding/json"
	"fmt"
	"github.com/ciazhar/go-start-small/examples/redpanda_schema_registry/internal/model"
	"github.com/ciazhar/go-start-small/examples/redpanda_schema_registry/proto"
	proto2 "google.golang.org/protobuf/proto"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/linkedin/goavro"
)

// Serialize JSON
func serializeJSON(user model.User) []byte {
	data, _ := json.Marshal(user)
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
	data, _ := proto2.Marshal(userProto)
	return data
}

// Produce message to Kafka
func produceMessage(topic string, value []byte) {

	// Kafka Config
	var kafkaConfig = &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092,localhost:9093,localhost:9094",
	}

	producer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}

	err = producer.Produce(msg, nil)
	if err != nil {
		log.Fatal(err)
	}
	producer.Flush(1000)
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
