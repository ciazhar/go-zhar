package postgres

import (
	"context"
	"embed"

	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func InitDBMigration(url string, fs embed.FS) {

	logger.LogInfo(context.Background(), "Running migrations", map[string]interface{}{
		"url": url})

	d, err := iofs.New(fs, ".")
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to initialize migration", nil)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, url)
	if err != nil {
		logger.LogFatal(context.Background(), err, "Failed to initialize migration", nil)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.LogFatal(context.Background(), err, "Failed to run migrations", nil)
	}

	logger.LogInfo(context.Background(), "Migrations completed", nil)
}
