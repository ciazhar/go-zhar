package main

import (
	"context"
	"flag"

	"github.com/ciazhar/go-start-small/examples/jwt_authorization_service/db/migrations"
	"github.com/ciazhar/go-start-small/examples/jwt_authorization_service/internal"
	"github.com/ciazhar/go-start-small/pkg/config"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/middleware"
	"github.com/ciazhar/go-start-small/pkg/postgres"
	"github.com/ciazhar/go-start-small/pkg/redis"
	"github.com/ciazhar/go-start-small/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
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

	pg := postgres.InitPostgres(
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.dbname"),
		viper.GetString("postgres.username"),
		viper.GetString("postgres.password"),
		logLevel,
	)
	defer pg.Close()
	postgres.InitDBMigration(
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.dbname"),
		viper.GetString("postgres.username"),
		viper.GetString("postgres.password"),
		migrations.MigrationsFS,
	)
	r := redis.InitRedis(
		viper.GetString("redis.host"),
		viper.GetInt("redis.port"),
		viper.GetString("redis.password"),
	)
	defer r.Close()
	v := validator.New("id")

	app := fiber.New()
	app.Use(middleware.RequestIDMiddleware)
	v1 := app.Group("/api/v1")
	internal.Init(v1, pg, r, v)

	// Start the server
	port := viper.GetString("server.port")
	if port == "" {
		port = "3000"
	}
	err := app.Listen(":" + port)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Server stopped", nil)
	}
}
