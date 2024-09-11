package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ciazhar/go-zhar/examples/mysql/transactional/model"
	"github.com/ciazhar/go-zhar/examples/mysql/transactional/repository"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func TestMain(m *testing.M) {
	ctx := context.Background()
	mysqlC, dbConn := setupTestContainer(ctx)
	db = dbConn
	code := m.Run()

	mysqlC.Terminate(ctx) // Cleanup after tests
	db.Close()
	os.Exit(code)
}

func setupTestContainer(ctx context.Context) (testcontainers.Container, *sql.DB) {
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      "testdb",
		},
		WaitingFor: wait.ForListeningPort("3306/tcp").WithStartupTimeout(2 * time.Minute),
	}

	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}

	port, _ := mysqlC.MappedPort(ctx, "3306")
	hostIP, _ := mysqlC.Host(ctx)

	dsn := fmt.Sprintf("root:root@tcp(%s:%s)/testdb?parseTime=true", hostIP, port.Port())
	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not connect to MySQL: %s", err)
	}

	return mysqlC, dbConn
}

func clearDatabase() {
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		log.Fatalf("Failed to clear database: %s", err)
	}
}

func TestCreateUserTransactionalSuccess(t *testing.T) {

	// Create the users table
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id BIGINT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) NOT NULL
    );`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Could not create users table: %s", err)
	}

	// Use transaction for repository
	repo := repository.NewMySQLUserRepository(db)

	user := &model.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	repo.WithTransaction(func(tx *sql.Tx) error {
		id, err := repo.Create(user)
		if err != nil {
			t.Fatalf("Failed to create user: %s", err)
		}

		if id == 0 {
			t.Fatalf("Expected a valid ID, got %d", id)
		}

		id2, err := repo.Create(user)
		if err != nil {
			t.Fatalf("Expected an error, got nil")
		}

		if id2 == 1 {
			t.Fatalf("Expected a valid ID, got %d", id)
		}

		return nil
	})

}

func TestCreateUserTransactionalError(t *testing.T) {

	// Create the users table
	dropTableQuery := `DROP TABLE users;`

	_, err := db.Exec(dropTableQuery)
	if err != nil {
		log.Fatalf("Could not create users table: %s", err)
	}

	// Create the users table
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id BIGINT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) NOT NULL UNIQUE
    );`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Could not create users table: %s", err)
	}

	// Use transaction for repository
	repo := repository.NewMySQLUserRepository(db)

	user := &model.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	repo.WithTransaction(func(tx *sql.Tx) error {
		id, err := repo.Create(user)
		if err != nil {
			t.Fatalf("Failed to create user: %s", err)
		}

		if id == 0 {
			t.Fatalf("Expected a valid ID, got %d", id)
		}

		_, err = repo.Create(user)
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}

		return nil
	})

}
