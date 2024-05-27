package main

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/internal/order"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/logger"
	opentelemetryjaeger "github.com/ciazhar/go-zhar/pkg/opentelemetry-jaeger"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func main() {
	log := logger.Init(logger.Config{
		ConsoleLoggingEnabled: true,
	})

	env.Init("config-order.json", log)

	tp, tracer := opentelemetryjaeger.InitTracer(
		viper.GetString("application.name"),
		viper.GetString("opentelemetry.url"),
		log,
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Infof("Error shutting down tracer provider: %v", err)
		}
	}()

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

	order.Init(app, client, tracer)

	err := app.Listen(":" + viper.GetString("application.port"))
	if err != nil {
		log.Fatalf("fiber failed to start : %v", err)
	}
}
