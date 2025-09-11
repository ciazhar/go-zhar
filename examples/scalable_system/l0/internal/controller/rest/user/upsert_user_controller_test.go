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

func TestUpsertUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockuser.NewMockUserService(ctrl)
	controller := NewUserController(mockService)

	app := fiber.New()
	app.Post("/users/upsert", func(c *fiber.Ctx) error {
		var body request.UpsertUserBodyRequest
		if err := c.BodyParser(&body); err != nil {
			return err
		}
		c.Locals("body", body)
		return controller.UpsertUser(c)
	})

	tests := []struct {
		name       string
		body       request.UpsertUserBodyRequest
		mockFunc   func()
		wantStatus int
		wantBody   string
	}{
		{
			name: "success upsert user",
			body: request.UpsertUserBodyRequest{
				Id:       "user-1",
				FullName: "Alice",
				Email:    "alice@example.com",
			},
			mockFunc: func() {
				mockService.EXPECT().
					UpsertUserByID(gomock.Any(), request.UpsertUserBodyRequest{
						Id:       "user-1",
						FullName: "Alice",
						Email:    "alice@example.com",
					}).
					Return(nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   "Upsert user success",
		},
		{
			name: "failed upsert user",
			body: request.UpsertUserBodyRequest{
				Id:       "user-2",
				FullName: "Bob",
				Email:    "bob@example.com",
			},
			mockFunc: func() {
				mockService.EXPECT().
					UpsertUserByID(gomock.Any(), request.UpsertUserBodyRequest{
						Id:       "user-2",
						FullName: "Bob",
						Email:    "bob@example.com",
					}).
					Return(errors.New("upsert failed"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "failed to upsert user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/users/upsert", bytes.NewBuffer(bodyBytes))
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
