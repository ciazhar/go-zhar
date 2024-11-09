package main

import (
	"github.com/ciazhar/go-start-small/examples/mongodb_migration/db/migrations"
	mongo "github.com/ciazhar/go-start-small/pkg/mongodb"
)

func main() {
	// Initialize database
	mongo.InitDBMigration("127.0.0.1:27017,127.0.0.1:27018,127.0.0.1:27019", "location", "", "", migrations.MigrationsFS)
}
