package user

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	mockuser "github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/service/user/mock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserController_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockuser.NewMockUserService(ctrl)
	controller := NewUserController(mockService)

	app := fiber.New()

	// Middleware to inject body into ctx.Locals
	app.Post("/users", func(c *fiber.Ctx) error {
		body := request.CreateUserBodyRequest{
			Username: "hafidz",
			Email:    "hafidz@example.com",
			Password: "hashed-password",
			FullName: "Muhammad Hafidz",
		}
		c.Locals("body", body)
		return controller.CreateUser(c)
	})

	tests := []struct {
		name       string
		mockFunc   func(mockService *mockuser.MockUserService)
		wantStatus int
		wantBody   string
	}{
		{
			name: "success create user",
			mockFunc: func(mockService *mockuser.MockUserService) {
				mockService.
					EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantStatus: http.StatusCreated,
			wantBody:   `"message":"Create user success"`,
		},
		{
			name: "failed create user",
			mockFunc: func(mockService *mockuser.MockUserService) {
				mockService.
					EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(errors.New("db insert error"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `"message":"failed to insert user to DB"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mockService)

			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{}`))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			bodyBytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			body := string(bodyBytes)
			assert.Contains(t, body, tt.wantBody)
		})
	}
}
