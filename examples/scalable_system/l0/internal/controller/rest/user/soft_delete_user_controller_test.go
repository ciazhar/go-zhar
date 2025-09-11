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
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockuser.NewMockUserService(ctrl)
	controller := NewUserController(mockService)

	app := fiber.New()
	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		// inject path param as ctx locals
		c.Locals("path_param", request.UserPathParam{ID: c.Params("id")})
		return controller.DeleteUser(c)
	})

	tests := []struct {
		name       string
		userID     string
		mockFunc   func()
		wantStatus int
		wantBody   string
	}{
		{
			name:   "success delete user",
			userID: "user-1",
			mockFunc: func() {
				mockService.EXPECT().
					SoftDeleteUser(gomock.Any(), "user-1").
					Return(nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   "Delete user success",
		},
		{
			name:   "failed delete user",
			userID: "user-2",
			mockFunc: func() {
				mockService.EXPECT().
					SoftDeleteUser(gomock.Any(), "user-2").
					Return(errors.New("delete failed"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "failed to delete user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.userID, bytes.NewBuffer(nil))
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
