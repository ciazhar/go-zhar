package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap/zapcore"
	"log"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

func initTracer() (func(context.Context) error, error) {
	// Set up debug logging
	logger, _ := zap.NewDevelopment(zap.IncreaseLevel(zapcore.DebugLevel))
	defer logger.Sync()

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		logger.Fatal("Failed to create Jaeger exporter", zap.Error(err))
		return nil, fmt.Errorf("creating Jaeger exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("your-service-name"),
		)),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Ensure all traces are sampled
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

func connectDatabase(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	// Configure the connection to use OpenTelemetry tracing
	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	conn, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return conn, nil
}

func main() {
	// Initialize the tracer
	shutdown, err := initTracer()
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer shutdown(context.Background())

	// Create a context
	ctx := context.Background()

	// Database connection string
	connString := "postgres://postgres:postgres@localhost:5432/postgres"

	// Connect to the database
	conn, err := connectDatabase(ctx, connString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close()

	// Simulate a database operation to generate traces
	if _, err := conn.Exec(ctx, "SELECT 1"); err != nil {
		log.Fatalf("failed to execute query: %v", err)
	}

	// Your application logic here
}
