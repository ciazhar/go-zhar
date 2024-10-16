package postgres

import (
	"context"
	"embed"
	"fmt"

	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func InitDBMigration(host string, port int, dbName string, username string, password string, fs embed.FS) {

	// Connection string
	connString := ""
	if username == "" && password == "" {
		connString = fmt.Sprintf("postgresql://%s:%d/%s?sslmode=disable", host, port, dbName)
	} else {
		connString = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, dbName)
	}

	logger.LogInfo(context.Background(), "Running migrations", map[string]interface{}{
		"url": connString})

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
