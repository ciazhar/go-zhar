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

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockuser.NewMockUserService(ctrl)
	controller := NewUserController(mockService)

	app := fiber.New()
	app.Put("/users/:id", func(c *fiber.Ctx) error {
		// inject locals like middleware
		c.Locals("path_param", request.UserPathParam{ID: c.Params("id")})

		var body request.UpdateUserBodyRequest
		if err := c.BodyParser(&body); err != nil {
			return err
		}
		c.Locals("body", body)

		return controller.UpdateUser(c)
	})

	tests := []struct {
		name       string
		userID     string
		body       request.UpdateUserBodyRequest
		mockFunc   func()
		wantStatus int
		wantBody   string
	}{
		{
			name:   "success update user",
			userID: "user-1",
			body: request.UpdateUserBodyRequest{
				FullName: "Alice",
				Email:    "alice@example.com",
			},
			mockFunc: func() {
				mockService.EXPECT().
					UpdateUser(gomock.Any(), "user-1", request.UpdateUserBodyRequest{
						FullName: "Alice",
						Email:    "alice@example.com",
					}).
					Return(nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   "Update user success",
		},
		{
			name:   "failed update user",
			userID: "user-2",
			body: request.UpdateUserBodyRequest{
				FullName: "Bob",
				Email:    "bob@example.com",
			},
			mockFunc: func() {
				mockService.EXPECT().
					UpdateUser(gomock.Any(), "user-2", request.UpdateUserBodyRequest{
						FullName: "Bob",
						Email:    "bob@example.com",
					}).
					Return(errors.New("update failed"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "failed to update user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/users/"+tt.userID, bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

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
