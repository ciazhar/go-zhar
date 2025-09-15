package jaeger

import (
	"context"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// InitJaegerTracer sets up Jaeger exporter
func InitJaegerTracer(ctx context.Context, applicationName, version, jaegerEndpoint string) func() {
	var (
		log = logger.FromContext(ctx)
	)

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(jaegerEndpoint),
	))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize Jaeger exporter")
	}

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(applicationName),
			semconv.ServiceVersion(version),
		),
	)
	if err != nil {
		panic(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return func() {
		_ = tp.Shutdown(ctx)
	}
}
