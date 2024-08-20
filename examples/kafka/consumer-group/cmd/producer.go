package main

import (
	"encoding/json"
	"github.com/ciazhar/go-zhar/examples/kafka/consumer-group/internal/event/model"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/kafka"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/spf13/viper"
	"time"
)

func main() {

	// Logger
	log := logger.Init(logger.Config{
		ConsoleLoggingEnabled: true,
	})

	// Environment configuration
	env.Init("config.json", log)

	producer := kafka.NewSyncProducer(viper.GetString("kafka.brokers"), log)
	defer producer.Close()

	event := model.EmailEvent{
		AmpEnabled:            true,
		BounceClass:           "hard",
		CampaignID:            "camp12345",
		ClickTracking:         true,
		CustomerID:            "cust67890",
		DelvMethod:            "SMTP",
		DeviceToken:           "device-token-123",
		ErrorCode:             "550",
		EventID:               "event12345",
		FriendlyFrom:          "no-reply@example.com",
		InitialPixel:          true,
		InjectionTime:         time.Now(),
		IPAddress:             "192.168.1.1",
		IPPool:                "pool1",
		MailboxProvider:       "Gmail",
		MailboxProviderRegion: "NA",
		MessageID:             "msg123456789",
		MsgFrom:               "sender@example.com",
		MsgSize:               "1048576", // 1 MB
		NumRetries:            "3",
		OpenTracking:          true,
		RCPTMeta:              map[string]interface{}{"key": "value"},
		RCPTTags:              []string{"tag1", "tag2"},
		RCPTTo:                "recipient@example.com",
		RCPTHash:              "rcpt-hash-123",
		RawRCPTTo:             "recipient@example.com",
		RCPTType:              "primary",
		RawReason:             "Mailbox full",
		Reason:                "bounce",
		RecipientDomain:       "example.com",
		RecvMethod:            "POP3",
		RoutingDomain:         "routing.example.com",
		ScheduledTime:         "2024-08-20T10:00:00Z",
		SendingDomain:         "example.com",
		SendingIP:             "203.0.113.1",
		SMSCoding:             "GSM-7",
		SMSDst:                "1234567890",
		SMSDstNPI:             "1",
		SMSDstTON:             "1",
		SMSSrc:                "0987654321",
		SMSSrcNPI:             "1",
		SMSSrcTON:             "1",
		SubaccountID:          "subacc123",
		Subject:               "Hello World",
		TemplateID:            "tmpl123",
		TemplateVersion:       "1.0",
		Timestamp:             time.Now(),
		Transactional:         "true",
		TransmissionID:        "trans123456",
		Type:                  "delivered",
	}

	marshal, err := json.Marshal(event)
	if err != nil {
		return
	}

	producer.PublishMessage("event", string(marshal))
}
