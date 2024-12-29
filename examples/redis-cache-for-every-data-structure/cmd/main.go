package main

import (
	"flag"
	"github.com/ciazhar/go-start-small/examples/redis-cache-for-every-data-structure/internal/controller"
	"github.com/ciazhar/go-start-small/examples/redis-cache-for-every-data-structure/internal/repository"
	"github.com/ciazhar/go-start-small/pkg/config"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"log"
)

func main() {

	// Configuration using flags for source, type, and other details
	var logLevel string
	var consoleOutput bool
	var source, configType, fileName, filePath, consulEndpoint, consulPath string

	// Parse command-line flags
	flag.StringVar(&logLevel, "log-level", "debug", "Log level (default: debug)")
	flag.BoolVar(&consoleOutput, "console-output", true, "Console output (default: true)")
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

	// Redis configuration
	r := redis.InitRedis(
		viper.GetString("redis.host"),
		viper.GetInt("redis.port"),
		viper.GetString("redis.password"),
	)
	defer r.Close()

	// Fiber configuration
	app := fiber.New()

	// Module initialization
	bitmapRepo := repository.NewBitmapDailyLoginTrackingRepository(r)
	geoRepo := repository.NewGeoTrackingDeliveryAgentLocationRepository(r)
	hashRepo := repository.NewHashUserProfileRepository(r)
	listRepo := repository.NewListOrderQueueRepository(r)
	setRepo := repository.NewSetActiveDeliveryAgentRepository(r)
	sortedSetRepo := repository.NewSortedSetRealtimeDeliveryMetricRepository(r)
	stringRepo := repository.NewStringUserSessionRepository(r)

	controller.NewDailyLoginTrackingController(bitmapRepo).RegisterRoutes(app)
	controller.NewGeoTrackingController(geoRepo).RegisterRoutes(app)
	controller.NewUserProfileController(hashRepo).RegisterRoutes(app)
	controller.NewOrderQueueController(listRepo).RegisterRoutes(app)
	controller.NewDeliveryAgentController(setRepo).RegisterRoutes(app)
	controller.NewDeliveryMetricController(sortedSetRepo).RegisterRoutes(app)
	controller.NewUserSessionController(stringRepo).RegisterRoutes(app)

	// Start Fiber
	err := app.Listen(":" + viper.GetString("application.port"))
	if err != nil {
		log.Fatalf("fiber failed to start : %v", err)
	}
}
