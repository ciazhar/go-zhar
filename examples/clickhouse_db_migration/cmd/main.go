package main

import (
	"github.com/ciazhar/go-start-small/examples/clickhouse_db_migration/db/migrations"
	"github.com/ciazhar/go-start-small/pkg/clickhouse"
)

func main() {
	// Initialize database
	clickhouse.InitDBMigration("localhost", 9000, "default", "", "", migrations.MigrationsFS)
}
