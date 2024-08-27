package repository

import (
	"errors"
	"strings"
)

// User represents a user with an ID, username, and password.
type User struct {
	ID       int
	Username string
	Password string
}

// AuthRepository is an interface for authentication-related operations.
type AuthRepository interface {
	FindByUsername(username string) (*User, error)
}

// InMemoryAuthRepository is an in-memory implementation of AuthRepository.
type InMemoryAuthRepository struct {
	users []User
}

// NewInMemoryAuthRepository creates a new instance of InMemoryAuthRepository with predefined users.
func NewInMemoryAuthRepository() *InMemoryAuthRepository {
	return &InMemoryAuthRepository{
		users: []User{
			{ID: 1, Username: "admin", Password: "password"},
			{ID: 2, Username: "user", Password: "1234"},
		},
	}
}

// FindByUsername searches for a user by their username.
func (repo *InMemoryAuthRepository) FindByUsername(username string) (*User, error) {
	for _, user := range repo.users {
		if strings.EqualFold(user.Username, username) {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}
