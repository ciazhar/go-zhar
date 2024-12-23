package main

import (
	"github.com/ciazhar/go-start-small/examples/audit_trail/internal/model"
	"github.com/ciazhar/go-start-small/pkg/kafka/producer"
	"log"
	"strings"
)

func main() {
	maxRetry := 5
	brokers := "localhost:8097,localhost:8098,localhost:8099"
	syncProducer, err := producer.NewSyncProducer(strings.Split(brokers, ","), maxRetry)
	if err != nil {
		log.Printf("Failed to create producer: %v", err)
		return
	}

	data := model.AuditLog{
		EventType: "user_login",
		UserID:    "1234",
		Timestamp: "2024-12-20T10:00:00Z",
	}
	_, _, err = syncProducer.SendMessage(model.AuditLogsTopic, "", data)
	if err != nil {
		return
	}

	log.Printf("Message sent successfully")

}
