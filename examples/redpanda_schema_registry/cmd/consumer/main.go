package main

import (
	"encoding/json"
	"fmt"
	"github.com/ciazhar/go-start-small/examples/redpanda_schema_registry/internal/model"
	"github.com/ciazhar/go-start-small/examples/redpanda_schema_registry/proto"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/linkedin/goavro"
	proto2 "google.golang.org/protobuf/proto"
	"log"
)

// Deserialize JSON
func deserializeJSON(data []byte) (model.User, error) {
	var user model.User
	err := json.Unmarshal(data, &user)
	return user, err
}

// Deserialize Avro
func deserializeAvro(data []byte) (model.User, error) {
	codec, err := goavro.NewCodec(model.AvroSchema)
	if err != nil {
		return model.User{}, err
	}
	native, _, err := codec.NativeFromBinary(data)
	if err != nil {
		return model.User{}, err
	}

	nativeMap := native.(map[string]interface{})
	return model.User{
		ID:    int(nativeMap["id"].(int32)),
		Name:  nativeMap["name"].(string),
		Email: nativeMap["email"].(string),
	}, nil
}

// Deserialize Protobuf
func deserializeProtobuf(data []byte) (model.User, error) {
	var userProto proto.UserProto
	err := proto2.Unmarshal(data, &userProto)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:    int(userProto.Id),
		Name:  userProto.Name,
		Email: userProto.Email,
	}, nil
}

// Kafka Consumer
func consume(topic string, deserializeFunc func([]byte) (model.User, error)) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092,localhost:9093,localhost:9094",
		"group.id":          "benchmark-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			user, err := deserializeFunc(msg.Value)
			if err != nil {
				log.Printf("Failed to deserialize message from topic %s: %v", topic, err)
				continue
			}
			fmt.Printf("Topic: %s | Size: %d bytes | User: %+v\n", topic, len(msg.Value), user)
			break
		}
	}
}

func main() {
	consume("benchmark-json", deserializeJSON)
	consume("benchmark-avro", deserializeAvro)
	consume("benchmark-protobuf", deserializeProtobuf)
}
