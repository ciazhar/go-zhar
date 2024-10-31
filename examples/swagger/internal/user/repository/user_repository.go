package repository

import (
	"fmt"
	"github.com/ciazhar/go-start-small/examples/swagger/internal/user/model"
)

// UserRepository represents a repository for managing users
type UserRepository struct {
	users map[string]model.User
}

// AddUser adds a new user to the repository
func (r *UserRepository) AddUser(user model.User) {
	r.users[user.Username] = user
}

// GetUserByUsername retrieves a user by their username
func (r *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// GetAllUsers retrieves all users from the repository
func (r *UserRepository) GetAllUsers() map[string]model.User {
	return r.users
}

// DeleteUser deletes a user from the repository
func (r *UserRepository) DeleteUser(username string) error {
	_, ok := r.users[username]
	if !ok {
		return fmt.Errorf("user not found for username: %s", username)
	}

	delete(r.users, username)
	return nil
}

func (r *UserRepository) UpdateUser(user model.User) error {
	r.users[user.Username] = user
	return nil
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]model.User),
	}
}
