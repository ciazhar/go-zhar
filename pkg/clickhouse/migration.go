package clickhouse

import (
	"context"
	"embed"
	"fmt"

	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// InitDBMigration initializes ClickHouse database migrations
func InitDBMigration(host string, port int, dbName string, username string, password string, fs embed.FS) {
	// Connection string for ClickHouse
	connString := ""
	if username == "" && password == "" {
		connString = fmt.Sprintf("clickhouse://%s:%d/%s", host, port, dbName)
	} else {
		connString = fmt.Sprintf("clickhouse://%s:%s@%s:%d/%s", username, password, host, port, dbName)
	}

	logger.LogInfo(context.Background(), "Running migrations", map[string]interface{}{
		"url": connString,
	})

	d, err := iofs.New(fs, ".")
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to initialize migration", nil)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, connString)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to initialize migration", nil)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.LogFatal(context.Background(), err, "Failed to run migrations", nil)
	}

	logger.LogInfo(context.Background(), "Migrations completed", nil)
}
