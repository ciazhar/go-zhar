package main

import (
	"context"
	"net/http"
	"time"

	"github.com/ciazhar/go-start-small/examples/http_distributed_tracing/internal/order"
	"github.com/ciazhar/go-start-small/pkg/config"
	"github.com/ciazhar/go-start-small/pkg/logger"
	oteljaeger "github.com/ciazhar/go-start-small/pkg/otel_jaeger"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {

	config.InitConfig(
		config.Config{
			Source: "file",
			Type:   "json",
			File: config.FileConfig{
				FileName: "config-order.json",
				FilePath: "configs",
			},
		},
	)

	traceProvider := oteljaeger.StartTracing("order-service", viper.GetString("opentelemetry.url"))
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			logger.LogFatal(context.Background(), err, "failed to shutdown", nil)
		}
	}()

	_ = traceProvider.Tracer("order-service")

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,               // Maximum number of idle connections
			IdleConnTimeout:     60 * time.Second, // Maximum time an idle connection is kept alive
			TLSHandshakeTimeout: 10 * time.Second, // Maximum time to wait for a TLS handshake
			DisableCompression:  true,             // Disable HTTP compression to improve performance
		},
	}

	app := fiber.New()

	order.Init(app, client)

	err := app.Listen(":" + viper.GetString("application.port"))
	if err != nil {
		logger.LogFatal(context.Background(), err, "fiber failed to start", nil)
	}
}
