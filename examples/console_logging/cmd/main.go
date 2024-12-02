package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/ciazhar/go-start-small/pkg/context_util"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/google/uuid"
)

func main() {
	// Initialize the logger
	logConfig := logger.LogConfig{
		ConsoleOutput: true,
		LogLevel:      "debug",
	}
	logger.InitLogger(logConfig)

	// Simulate a login request with a unique request ID
	ctx := context.WithValue(context.Background(), context_util.RequestIDKey, uuid.New().String())

	// Mock user credentials
	username := "testuser"
	password := "password123"

	// Call the login function
	if err := login(ctx, username, password); err != nil {
		fmt.Println("Login process completed with errors.")
	} else {
		fmt.Println("Login successful.")
	}
}

// Mock database of users
var mockUserDB = map[string]string{
	"testuser": "password123",
}

// login simulates a login process with logging for different scenarios
func login(ctx context.Context, username, password string) error {
	// Log the start of the login process
	logger.LogInfo(ctx, "Starting login process", map[string]interface{}{"username": username})

	// Simulate user lookup
	storedPassword, userExists := mockUserDB[username]
	if !userExists {
		// Log and return a warning for unknown username
		return logger.LogAndReturnWarning(ctx, errors.New("user not found"), "Invalid login attempt: user does not exist", map[string]interface{}{"username": username})
	}

	// Simulate password validation
	if password != storedPassword {
		// Log and return a warning for incorrect password
		return logger.LogAndReturnWarning(ctx, errors.New("invalid credentials"), "Invalid login attempt: incorrect password", map[string]interface{}{"username": username})
	}

	// Log successful login
	logger.LogInfo(ctx, "User successfully logged in", map[string]interface{}{"username": username})

	return nil
}
