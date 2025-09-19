package user

import (
	"bytes"
	"encoding/json"
	"errors"
	mockuser "github.com/ciazhar/go-zhar-scalable-system-l0/internal/service/user/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ciazhar/go-zhar-scalable-system-l0/internal/model/request"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestIsUserExistByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockuser.NewMockUserService(ctrl)
	controller := NewUserController(mockService)

	app := fiber.New()
	app.Get("/users/exist", func(c *fiber.Ctx) error {
		// inject query param as ctx locals
		c.Locals("query_param", request.UserEmailQueryParam{Email: c.Query("email")})
		return controller.IsUserExistByEmail(c)
	})

	tests := []struct {
		name       string
		email      string
		mockFunc   func()
		wantStatus int
		wantBody   string
	}{
		{
			name:  "user exists",
			email: "hafidz@example.com",
			mockFunc: func() {
				mockService.EXPECT().
					IsUserExistsByEmail(gomock.Any(), "hafidz@example.com").
					Return(true, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   "Check user existence success",
		},
		{
			name:  "user does not exist",
			email: "unknown@example.com",
			mockFunc: func() {
				mockService.EXPECT().
					IsUserExistsByEmail(gomock.Any(), "unknown@example.com").
					Return(false, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   "Check user existence success",
		},
		{
			name:  "service error",
			email: "error@example.com",
			mockFunc: func() {
				mockService.EXPECT().
					IsUserExistsByEmail(gomock.Any(), "error@example.com").
					Return(false, errors.New("db error"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "failed to check user existence",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			req := httptest.NewRequest(http.MethodGet, "/users/exist?email="+tt.email, bytes.NewBuffer(nil))
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
