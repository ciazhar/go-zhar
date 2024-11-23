package main

import (
	"flag"
	"github.com/ciazhar/go-start-small/examples/kafka_sync_async_producer/internal"
	"github.com/ciazhar/go-start-small/examples/kafka_sync_async_producer/internal/model"
	"github.com/ciazhar/go-start-small/pkg/config"
	"github.com/ciazhar/go-start-small/pkg/kafka"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"log"
	"strings"
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

	// Kafka configuration
	admin := kafka.CreateKafkaAdminClient(strings.Split(viper.GetString("kafka.brokers"), ","))
	kafka.CreateKafkaTopic(admin, model.TopicSync, 3, 1, 0, nil)
	kafka.CreateKafkaTopic(admin, model.TopicAsync, 3, 1, 0, nil)

	syncProducer := kafka.CreateProducer(strings.Split(viper.GetString("kafka.brokers"), ","))
	asyncProducer := kafka.CreateAsyncProducer(strings.Split(viper.GetString("kafka.brokers"), ","))

	app := fiber.New()
	internal.Init(app, syncProducer, asyncProducer)

	err := app.Listen(":" + viper.GetString("application.port"))
	if err != nil {
		log.Fatalf("fiber failed to start : %v", err)
	}
}
