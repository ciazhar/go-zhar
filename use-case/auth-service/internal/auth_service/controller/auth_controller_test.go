package controller_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/ciazhar/go-zhar/use-case/auth-service/internal/auth_service/controller"
	"github.com/ciazhar/go-zhar/use-case/auth-service/internal/auth_service/model"
	"github.com/ciazhar/go-zhar/use-case/auth-service/internal/mocks"
	"github.com/ciazhar/go-zhar/use-case/auth-service/pkg/middleware"
	"github.com/ciazhar/go-zhar/use-case/auth-service/pkg/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gofiber/fiber/v2"
)

// Test RegisterUser
func TestRegisterUser(t *testing.T) {

	validation.InitValidation()

	app := fiber.New()
	app.Use(middleware.RequestIDMiddleware)

	mockAuthService := new(mocks.AuthServiceInterface)
	authController := controller.NewAuthController(mockAuthService)

	app.Post("/register", authController.RegisterUser)

	// Define test cases
	tests := []struct {
		name       string
		body       []byte
		statusCode int
		mockFunc   func()
	}{
		{
			name:       "Success",
			body:       []byte(`{"username": "testuser", "password": "password"}`),
			statusCode: fiber.StatusCreated,
			mockFunc: func() {
				mockAuthService.On("RegisterUser", mock.Anything, mock.AnythingOfType("model.User")).
					Return(nil)
			},
		},
		{
			name:       "Invalid Input",
			body:       []byte(`{"username": "", "password": "password"}`),
			statusCode: fiber.StatusBadRequest,
			mockFunc:   func() {},
		},
		{
			name:       "Error on Register",
			body:       []byte(`{"username": "testuser", "password": "password"}`),
			statusCode: fiber.StatusInternalServerError,
			mockFunc: func() {
				mockAuthService.On("RegisterUser", mock.Anything, mock.AnythingOfType("model.User")).
					Return(assert.AnError) // Simulating an error
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the mock function to set expectations
			tt.mockFunc()

			// Create a new request using httptest
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(tt.body))
			req.Header.Set("Content-Type", "application/json")

			// Perform the request
			resp, err := app.Test(req)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, tt.statusCode, resp.StatusCode)

			// Assert that the mock expectations were met
			mockAuthService.AssertExpectations(t)
		})
	}
}

// Test Login
func TestLogin(t *testing.T) {

	validation.InitValidation()

	app := fiber.New()

	mockAuthService := new(mocks.AuthServiceInterface)
	authController := controller.NewAuthController(mockAuthService)

	app.Post("/login", authController.Login)

	// Define test cases
	tests := []struct {
		name       string
		body       []byte
		statusCode int
		mockFunc   func()
	}{
		{
			name:       "Success",
			body:       []byte(`{"username": "testuser", "password": "password"}`),
			statusCode: fiber.StatusOK,
			mockFunc: func() {
				mockAuthService.On("Login", mock.Anything, mock.AnythingOfType("model.LoginRequest")).
					Return(model.LoginResponse{AccessToken: "some_access_token", RefreshToken: "some_refresh_token"}, nil)
			},
		},
		{
			name:       "Invalid Input",
			body:       []byte(`{"username": "", "password": "password"}`),
			statusCode: fiber.StatusBadRequest,
			mockFunc:   func() {},
		},
		{
			name:       "Error on Login",
			body:       []byte(`{"username": "testuser", "password": "password"}`),
			statusCode: fiber.StatusInternalServerError,
			mockFunc: func() {
				mockAuthService.On("Login", mock.Anything, mock.AnythingOfType("model.LoginRequest")).
					Return(model.LoginResponse{}, assert.AnError) // Simulating an error
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Create a new request using httptest
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(tt.body))
			req.Header.Set("Content-Type", "application/json")

			// Perform the request
			resp, err := app.Test(req)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, tt.statusCode, resp.StatusCode)
		})
	}
}
