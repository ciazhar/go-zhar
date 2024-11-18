package utils

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// InitDB initializes the database connection
func InitDB() *sql.DB {
	dsn := "postgres://@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Ensure the tasks table exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL UNIQUE,
			description TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	return db
}
