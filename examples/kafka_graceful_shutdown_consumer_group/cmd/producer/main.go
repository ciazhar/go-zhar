package main

import (
	"encoding/json"
	"flag"
	"github.com/ciazhar/go-start-small/examples/kafka_graceful_shutdown_consumer_group/internal/event/model"
	"github.com/ciazhar/go-start-small/pkg/config"
	"github.com/ciazhar/go-start-small/pkg/kafka"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func main() {

	// Configuration using flags for source, type, and other details
	var logLevel string
	var consoleOutput bool
	var source, configType, fileName, filePath, consulEndpoint, consulPath string

	// Parse command-line flags
	flag.StringVar(&logLevel, "log-level", "debug", "Log level")
	flag.BoolVar(&consoleOutput, "console-output", true, "Console output")
	flag.StringVar(&source, "source", "file", "Configuration source (file or consul)")
	flag.StringVar(&fileName, "file-name", "config.json", "Name of the configuration file")
	flag.StringVar(&filePath, "file-path", "configs", "Path to the configuration file")
	flag.StringVar(&configType, "config-type", "json", "Configuration file type")
	flag.StringVar(&consulEndpoint, "consul-endpoint", "localhost:8500", "Consul endpoint")
	flag.StringVar(&consulPath, "consul-path", "path/to/config", "Path to the configuration in Consul")
	flag.Parse()

	// Initialize logger with parsed configuration
	logger.InitLogger(logger.LogConfig{
		LogLevel:      logLevel,
		ConsoleOutput: consoleOutput,
	})

	// Configuration using flags for source, type, and other details
	fileConfig := config.Config{
		Source: source,
		Type:   configType,
		File: config.FileConfig{
			FileName: fileName,
			FilePath: filePath,
		},
		Consul: config.ConsulConfig{
			Endpoint: consulEndpoint,
			Path:     consulPath,
		},
	}

	config.InitConfig(fileConfig)

	syncProducer := kafka.CreateProducer(strings.Split(viper.GetString("kafka.brokers"), ","))
	defer syncProducer.Close()

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

	err = kafka.SendMessage(syncProducer, "my-topic", string(marshal))
	if err != nil {
		return
	}
}
