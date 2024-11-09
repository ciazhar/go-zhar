package main

import (
	"context"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ciazhar/go-start-small/examples/clickhouse_export_csv/internal/repository"
	"github.com/ciazhar/go-start-small/pkg/clickhouse"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/gofiber/fiber/v2/log"
)

func main() {

	// Logger
	logger.InitLogger(logger.LogConfig{
		ConsoleOutput: true,
	})

	conn := clickhouse.Init(
		"localhost:9000",
		"default",
		"default",
		"",
		true,
	)
	defer conn.Close()

	// Example usage
	ctx := context.Background()
	repo := &repository.ClickhouseRepository{
		conn: conn,
	}
	_, err := repo.ExportEvents(ctx, "", "")
	if err != nil {
		log.Fatalf("failed to export events: %s", err)
	}
}
