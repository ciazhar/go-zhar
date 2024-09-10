package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ciazhar/go-zhar/examples/mysql/crud/model"
	"github.com/ciazhar/go-zhar/examples/mysql/crud/repository"
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

	// Create the users table
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id BIGINT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) NOT NULL
    );`

	_, err = dbConn.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Could not create users table: %s", err)
	}

	return mysqlC, dbConn
}

func clearDatabase() {
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		log.Fatalf("Failed to clear database: %s", err)
	}
}

func TestCreateUser(t *testing.T) {
	clearDatabase()
	repo := repository.NewMySQLUserRepository(db)

	user := &model.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	id, err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}

	if id == 0 {
		t.Fatalf("Expected a valid ID, got %d", id)
	}
}

func TestGetUserByID(t *testing.T) {
	clearDatabase()
	repo := repository.NewMySQLUserRepository(db)

	user := &model.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	id, err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}

	user, err = repo.GetByID(id)
	if err != nil {
		t.Fatalf("Failed to get user: %s", err)
	}

	if user == nil {
		t.Fatalf("Expected to get a user, got nil")
	}

	if user.Name != "John Doe" {
		t.Errorf("Expected name to be 'John Doe', got '%s'", user.Name)
	}
}

func TestUpdateUser(t *testing.T) {
	clearDatabase()
	repo := repository.NewMySQLUserRepository(db)

	// Insert a user
	user := &model.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	id, err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}

	// Update the user's name
	updatedUser := &model.User{
		ID:    id,
		Name:  "Jane Doe",
		Email: "john@example.com", // Keeping email the same
	}

	err = repo.Update(updatedUser)
	if err != nil {
		t.Fatalf("Failed to update user: %s", err)
	}

	// Fetch the updated user
	fetchedUser, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("Failed to fetch updated user: %s", err)
	}

	if fetchedUser.Name != "Jane Doe" {
		t.Errorf("Expected user name to be 'Jane Doe', but got '%s'", fetchedUser.Name)
	}
}

func TestDeleteUser(t *testing.T) {
	clearDatabase()
	repo := repository.NewMySQLUserRepository(db)

	// Insert a user
	user := &model.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	id, err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}

	// Delete the user
	err = repo.Delete(id)
	if err != nil {
		t.Fatalf("Failed to delete user: %s", err)
	}

	// Attempt to fetch the deleted user
	_, err = repo.GetByID(id)
	if err == nil {
		t.Fatalf("Expected error when fetching deleted user, but got none")
	}

	if err != sql.ErrNoRows {
		t.Errorf("Expected sql.ErrNoRows, but got '%s'", err)
	}
}

func TestFindAll(t *testing.T) {
	clearDatabase()
	repo := repository.NewMySQLUserRepository(db)

	// Insert test data
	_, err := repo.Create(&model.User{Name: "John Doe", Email: "john@example.com"})
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}
	_, err = repo.Create(&model.User{Name: "Jane Doe", Email: "jane@example.com"})
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}

	nameFilter := "John"
	emailFilter := "john@example.com"
	limit := 10
	offset := 0

	users, err := repo.FindAll(&nameFilter, &emailFilter, limit, offset)
	if err != nil {
		t.Fatalf("Failed to find users: %s", err)
	}

	for i := range users {
		fmt.Printf("User: %s, Email: %s\n", users[i].Name, users[i].Email)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}
}

func TestCount(t *testing.T) {
	clearDatabase()
	repo := repository.NewMySQLUserRepository(db)

	// Insert test data
	_, err := repo.Create(&model.User{Name: "John Doe", Email: "john@example.com"})
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}
	_, err = repo.Create(&model.User{Name: "Jane Doe", Email: "jane@example.com"})
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}

	nameFilter := "Doe"
	emailFilter := "john"

	count, err := repo.Count(&nameFilter, &emailFilter)
	if err != nil {
		t.Fatalf("Failed to count users: %s", err)
	}

	if count != 1 {
		t.Errorf("Expected 2 users, got %d", count)
	}
}

func TestFindWithCursor(t *testing.T) {
	clearDatabase()
	repo := repository.NewMySQLUserRepository(db)

	// Insert test data
	id1, err := repo.Create(&model.User{Name: "John Doe", Email: "john@example.com"})
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}
	_, err = repo.Create(&model.User{Name: "Jane Doe", Email: "jane@example.com"})
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}

	limit := 10
	lastID := id1

	users, err := repo.FindWithCursor(&lastID, nil, nil, limit)
	if err != nil {
		t.Fatalf("Failed to find users with cursor: %s", err)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}
}
