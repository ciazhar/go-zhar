package user

import (
	"context"
	"errors"
	mockuser "github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/repository/postgres/user/mock"
	"testing"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockuser.NewMockUserRepositoryContract(ctrl)
	service := userService{repo: mockRepo}
	mockReq := request.CreateUserBodyRequest{
		Username: "hafidz",
		Email:    "hafidz@example.com",
		Password: "hashed-password",
		FullName: "Muhammad Hafidz",
	}

	tests := []struct {
		name     string
		input    request.CreateUserBodyRequest
		mockFunc func(mockRepo *mockuser.MockUserRepositoryContract)
		wantErr  bool
	}{
		{
			name:  "success insert user",
			input: mockReq,
			mockFunc: func(mockRepo *mockuser.MockUserRepositoryContract) {
				mockRepo.
					EXPECT().
					CreateUser(gomock.Any(), mockReq).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "failed insert user",
			input: mockReq,
			mockFunc: func(mockRepo *mockuser.MockUserRepositoryContract) {
				mockRepo.
					EXPECT().
					CreateUser(gomock.Any(), mockReq).
					Return(errors.New("db insert error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mockRepo)

			err := service.CreateUser(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
