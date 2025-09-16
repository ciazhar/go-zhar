package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ciazhar/go-zhar/pkg/jaeger"
	metrics "github.com/ciazhar/go-zhar/pkg/prometheus"
	"time"

	bootstrap2 "github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/bootstrap"
	"github.com/ciazhar/go-zhar/pkg/postgres"

	"github.com/ciazhar/go-zhar/pkg/bootstrap"
	"github.com/ciazhar/go-zhar/pkg/bootstrap/server"
	"github.com/ciazhar/go-zhar/pkg/config"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/validator"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()

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

	//Configure logger
	log := logger.FromContext(ctx)

	// Configuration using flags for source, type, and other details
	configSource, err := config.ParseConfigSource(source)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse config source")
	}
	fileConfig := config.Config{
		Source: configSource,
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
	config.InitConfig(ctx, fileConfig)

	// Initialize logger with log rotation
	logConfig := logger.LogConfig{
		LogLevel:      "debug",
		LogFile:       fmt.Sprintf("logs/%s.log", viper.GetString("application.name")),
		MaxSize:       1, // 1 MB
		MaxBackups:    3,
		MaxAge:        1,
		Compress:      true,
		ConsoleOutput: false,
	}
	logger.InitLogger(logConfig)

	// Initialize metric
	metrics.InitMetrics()

	// Initialize tracer
	shutdown := jaeger.InitJaegerTracer(ctx,
		viper.GetString("application.name"),
		viper.GetString("application.version"),
		viper.GetString("jaeger.endpoint"))
	defer shutdown()

	// === INIT VALIDATOR ===
	v := validator.New("id")

	// === INIT CLIENTS ===
	var clients []bootstrap.Service
	postgresSvc, postgresClient := postgres.InitPostgres(ctx,
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.dbname"),
		viper.GetString("postgres.username"),
		viper.GetString("postgres.password"),
		logLevel,
	)
	clients = append(clients, postgresSvc)

	// === START ALL CLIENTS ===
	for _, svc := range clients {
		if err := svc.Start(); err != nil {
			logger.LogFatal(err).Msgf("🔥 %s failed", svc.Name())
		}
	}

	// === BUILD HTTP MODULE WITH WIRE ===
	restModule := bootstrap2.InitializeRESTModule(v, postgresClient)

	// === INIT SERVERS + WORKERS ===
	var serversAndWorkers []bootstrap.Service
	server1 := server.NewFiberServer(ctx,
		fmt.Sprintf("%s:%s", viper.GetString("application.name"), viper.GetString("application.version")),
		fmt.Sprintf(":%s", viper.GetString("application.port")),
		restModule.Register, // <- just pass the registrar
	)
	serversAndWorkers = append(serversAndWorkers, server1)

	// === START ALL SERVERS + WORKERS ===
	for _, svc := range serversAndWorkers {
		go func(svc bootstrap.Service) {
			if err := svc.Start(); err != nil {
				log.Fatal().Msgf("🔥 %s failed", svc.Name())
			}
		}(svc)
	}

	// === GRACEFUL SHUTDOWN ===
	bootstrap.GracefulShutdown(ctx, 10*time.Second, clients, serversAndWorkers)
}
