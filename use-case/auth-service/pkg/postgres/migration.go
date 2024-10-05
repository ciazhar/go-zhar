package postgres

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

func InitPostgresMigration() {

	// Create a new migrate instance
	m, err := migrate.New("file://db/migrations", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	// Run the migrations
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No changes to apply in migrations.")
		} else {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	}

	log.Println("Migrations applied successfully!")
}
