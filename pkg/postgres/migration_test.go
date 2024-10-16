package postgres

// import (
// 	"context"
// 	"embed"
// 	"fmt"
// 	"log"
// 	"testing"
// 	"time"

// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/testcontainers/testcontainers-go"
// 	"github.com/testcontainers/testcontainers-go/wait"
// )

// //go:embed testdata/migrations/*.sql
// var testMigrations embed.FS

// func startPostgresContainer(t *testing.T) (testcontainers.Container, string, error) {
// 	ctx := context.Background()

// 	// Request a PostgreSQL container
// 	req := testcontainers.ContainerRequest{
// 		Image:        "postgres:16-alpine", // Use any PostgreSQL version
// 		ExposedPorts: []string{"5432/tcp"},
// 		Env: map[string]string{
// 			"POSTGRES_USER":     "testuser",
// 			"POSTGRES_PASSWORD": "testpassword",
// 			"POSTGRES_DB":       "testdb",
// 		},
// 		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(1),
// 	}

// 	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
// 		ContainerRequest: req,
// 		Started:          true,
// 	})
// 	if err != nil {
// 		return nil, "", fmt.Errorf("could not start container: %w", err)
// 	}

// 	host, err := container.Host(ctx)
// 	if err != nil {
// 		return container, "", err
// 	}

// 	port, err := container.MappedPort(ctx, "5432")
// 	if err != nil {
// 		return container, "", err
// 	}

// 	url := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", host, port.Port())
// 	return container, url, nil
// }

// func TestInitDBMigration(t *testing.T) {
// 	// Start PostgreSQL container
// 	postgresContainer, url, err := startPostgresContainer(t)
// 	if err != nil {
// 		log.Fatalf("Failed to start PostgreSQL container: %v", err)
// 	}
// 	defer func() {
// 		_ = postgresContainer.Terminate(context.Background()) // Ensure container is cleaned up after test
// 	}()

// 	// Embed the dynamically created migration file (as an alternative to embed.FS, we can just pass the file path)
// 	InitDBMigration(url, testMigrations) // Modify the function signature if necessary

// 	// Create a connection pool using pgxpool
// 	ctx := context.Background()
// 	poolConfig, err := pgxpool.ParseConfig(url)
// 	if err != nil {
// 		log.Fatalf("Failed to parse pgxpool config: %v", err)
// 	}

// 	// Set connection pool parameters
// 	poolConfig.MaxConns = 10
// 	poolConfig.MaxConnIdleTime = 5 * time.Second

// 	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
// 	if err != nil {
// 		log.Fatalf("Failed to create pgxpool connection pool: %v", err)
// 	}
// 	defer pool.Close()

// 	// Wait for the pool to be ready
// 	err = pool.Ping(ctx)
// 	if err != nil {
// 		log.Fatalf("Failed to ping the PostgreSQL database: %v", err)
// 	}

// 	// Query the database to verify if the migrations were applied
// 	var exists bool
// 	err = pool.QueryRow(ctx, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&exists)
// 	assert.NoError(t, err)
// 	assert.True(t, exists, "The migration should have created the 'users' table")
// }