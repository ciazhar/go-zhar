package main

import (
	"context"

	"github.com/ciazhar/go-start-small/examples/http_distributed_tracing/internal/user"
	"github.com/ciazhar/go-start-small/pkg/config"
	"github.com/ciazhar/go-start-small/pkg/logger"
	oteljaeger "github.com/ciazhar/go-start-small/pkg/otel_jaeger/v1"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {

	config.InitConfig(
		config.Config{
			Source: "file",
			Type:   "json",
			File: config.FileConfig{
				FileName: "config-user.json",
				FilePath: "configs",
			},
		},
	)

	traceProvider := oteljaeger.StartTracing("user-service", viper.GetString("opentelemetry.url"))
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			logger.LogFatal(context.Background(), err, "failed to shutdown", nil)
		}
	}()

	_ = traceProvider.Tracer("user-service")

	app := fiber.New()

	user.Init(app)

	err := app.Listen(":" + viper.GetString("application.port"))
	if err != nil {
		logger.LogFatal(context.Background(), err, "fiber failed to start", nil)
	}
}
