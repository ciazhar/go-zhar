package main

import (
	"github.com/ciazhar/go-start-small/examples/postgres_db_migration/db/migrations"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/postgres"
)

func main() {

	// Initialize logger
	testConfig := logger.LogConfig{
		ConsoleOutput: true,
		LogLevel:      "debug",
	}
	logger.InitLogger(testConfig)

	postgres.InitDBMigration("postgresql://localhost:5432/dbname?sslmode=disable", migrations.MigrationsFS)
}
