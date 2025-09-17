package user

import (
	"bytes"
	"encoding/json"
	"errors"
	mockuser "github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/service/user/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/response"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockuser.NewMockUserService(ctrl)
	controller := NewUserController(mockService)

	app := fiber.New()
	app.Get("/users", func(c *fiber.Ctx) error {
		// inject query param as ctx locals
		c.Locals("query_param", request.GetUsersQueryParam{Page: 1, Size: 10})
		return controller.GetUsers(c)
	})

	mockUsers := []response.User{
		{
			ID:       "user-1",
			Username: "hafidz",
			Email:    "hafidz@example.com",
			FullName: "Muhammad Hafidz",
		},
		{
			ID:       "user-2",
			Username: "vica",
			Email:    "vica@example.com",
			FullName: "Agnesia Vica",
		},
	}

	tests := []struct {
		name       string
		mockFunc   func()
		wantStatus int
		wantBody   string
	}{
		{
			name: "success get users",
			mockFunc: func() {
				mockService.EXPECT().
					GetUsersWithPagination(gomock.Any(), 1, 10).
					Return(mockUsers, int64(len(mockUsers)), nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   "Get users success",
		},
		{
			name: "failed get users",
			mockFunc: func() {
				mockService.EXPECT().
					GetUsersWithPagination(gomock.Any(), 1, 10).
					Return(nil, int64(0), errors.New("db error"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "failed to get users from DB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			req := httptest.NewRequest(http.MethodGet, "/users", bytes.NewBuffer(nil))
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			var resBody map[string]any
			err = json.NewDecoder(resp.Body).Decode(&resBody)
			assert.NoError(t, err)

			assert.Contains(t, resBody["message"], tt.wantBody)
		})
	}
}
