package service

import (
	"github.com/ciazhar/go-zhar/examples/fiber/clean-architecture-swagger/internal/user/model"
	"github.com/ciazhar/go-zhar/examples/fiber/clean-architecture-swagger/internal/user/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepo}
}

// AddUser adds a new user
func (s *UserService) AddUser(user model.User) {
	s.userRepository.AddUser(user)
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	return s.userRepository.GetUserByUsername(username)
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers() map[string]model.User {
	return s.userRepository.GetAllUsers()
}

// DeleteUser deletes a user by username
func (s *UserService) DeleteUser(username string) error {
	return s.userRepository.DeleteUser(username)
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(user model.User) error {
	return s.userRepository.UpdateUser(user)
}
