package user

import (
	"context"
	"errors"
	"testing"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	mock_user "github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/repository/postgres/user/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpsertUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRepositoryContract(ctrl)
	service := NewUserService(mockRepo)

	mockReq := request.UpsertUserBodyRequest{
		Id:       "user-1",
		Username: "hafidz",
		Email:    "hafidz@example.com",
		FullName: "Muhammad Hafidz",
	}

	tests := []struct {
		name     string
		input    request.UpsertUserBodyRequest
		mockFunc func(mockRepo *mock_user.MockUserRepositoryContract)
		wantErr  bool
	}{
		{
			name:  "success upsert user",
			input: mockReq,
			mockFunc: func(mockRepo *mock_user.MockUserRepositoryContract) {
				mockRepo.
					EXPECT().
					UpsertUserByID(gomock.Any(), mockReq).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "failed upsert user",
			input: mockReq,
			mockFunc: func(mockRepo *mock_user.MockUserRepositoryContract) {
				mockRepo.
					EXPECT().
					UpsertUserByID(gomock.Any(), mockReq).
					Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mockRepo)

			err := service.UpsertUserByID(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
