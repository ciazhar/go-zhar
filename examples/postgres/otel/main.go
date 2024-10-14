package main

import (
	"context"
	"fmt"
	"github.com/exaring/otelpgx"
	zerolog "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	zerolog2 "github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"os"
)

type MultiQueryTracer struct {
	Tracers []pgx.QueryTracer
}

func (m *MultiQueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	for _, t := range m.Tracers {
		ctx = t.TraceQueryStart(ctx, conn, data)
	}

	return ctx
}

func (m *MultiQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	for _, t := range m.Tracers {
		t.TraceQueryEnd(ctx, conn, data)
	}
}

func New(ctx context.Context, dsn string, log zerolog2.Logger) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}

	// FYI: NewLoggerAdapter: https://github.com/mcosta74/pgx-slog
	adapterLogger := zerolog.NewLogger(log)

	m := MultiQueryTracer{
		Tracers: []pgx.QueryTracer{
			// tracer: https://github.com/exaring/otelpgx
			otelpgx.NewTracer(),

			// logger
			&tracelog.TraceLog{
				Logger:   adapterLogger,
				LogLevel: tracelog.LogLevelTrace,
			},
		},
	}

	config.ConnConfig.Tracer = &m

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(err)
	}

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	return pool
}

func initTracer() (func(context.Context) error, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		return nil, fmt.Errorf("creating Jaeger exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

func main() {

	shutdown, err := initTracer()
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer shutdown(context.Background())

	pool := New(context.Background(), "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable", zerolog2.New(os.Stdout))

	// Run SQL
	_, err = pool.Exec(context.Background(), "SELECT 1")
	if err != nil {
		panic(err)
	}
}
