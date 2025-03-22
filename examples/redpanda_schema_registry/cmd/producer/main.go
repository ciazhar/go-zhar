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
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	brokers := []string{"localhost:9092", "localhost:9093", "localhost:9094"}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		log.Fatal(err)
	}
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
