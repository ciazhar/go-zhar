package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// User represents a user data structure.
type User struct {
	ID    int
	Name  string
	Email string
}

// UserPool manages a pool of User objects.
type UserPool struct {
	pool sync.Pool
}

// NewUserPool creates a new UserPool.
func NewUserPool() *UserPool {
	return &UserPool{
		pool: sync.Pool{
			New: func() interface{} {
				return new(User)
			},
		},
	}
}

// Get retrieves a User from the pool.
func (up *UserPool) Get() *User {
	return up.pool.Get().(*User)
}

// Put returns a User to the pool.
func (up *UserPool) Put(user *User) {
	user.ID = 0
	user.Name = ""
	user.Email = ""
	up.pool.Put(user)
}

// fetchAllUsersFromDB fetches all Users from the database.
func fetchAllUsersFromDB(db *sql.DB, up *UserPool) ([]*User, error) {
	query := "SELECT id, name, email FROM users"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0, 10) // Pre-allocate slice with initial capacity
	for rows.Next() {
		user := up.Get()
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			up.Put(user)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func main() {
	db, err := sql.Open("sqlite3", "example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table and insert some data (for demonstration purposes)
	createTable := `
 CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY,
  name TEXT,
  email TEXT
 )`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	insertData := `
 INSERT INTO users (id, name, email) VALUES (1, 'Alice', 'alice@example.com'), 
                                              (2, 'Bob', 'bob@example.com'), 
                                              (3, 'Charlie', 'charlie@example.com')
 ON CONFLICT(id) DO NOTHING
 `
	_, err = db.Exec(insertData)
	if err != nil {
		log.Fatal(err)
	}

	userPool := NewUserPool()

	// Fetch all Users from the database
	users, err := fetchAllUsersFromDB(db, userPool)
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		fmt.Printf("Fetched User with ID: %d, Name: %s, Email: %s\n", user.ID, user.Name, user.Email)
		userPool.Put(user) // Return the User to the pool
	}
}
