package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/response"
	mock_user "github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/repository/postgres/user/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRepositoryContract(ctrl)
	service := NewUserService(mockRepo)

	mockID := uuid.New().String()
	mockUser := &response.User{
		ID:       mockID,
		Username: "hafidz",
		Email:    "hafidz@example.com",
		FullName: "Muhammad Hafidz",
	}

	tests := []struct {
		name     string
		id       string
		mockFunc func(mockRepo *mock_user.MockUserRepositoryContract)
		wantUser *response.User
		wantErr  bool
	}{
		{
			name: "success get user by ID",
			id:   mockID,
			mockFunc: func(mockRepo *mock_user.MockUserRepositoryContract) {
				mockRepo.
					EXPECT().
					GetUserByID(gomock.Any(), mockID).
					Return(mockUser, nil)
			},
			wantUser: mockUser,
			wantErr:  false,
		},
		{
			name: "failed get user by ID",
			id:   mockID,
			mockFunc: func(mockRepo *mock_user.MockUserRepositoryContract) {
				mockRepo.
					EXPECT().
					GetUserByID(gomock.Any(), mockID).
					Return(nil, errors.New("db error"))
			},
			wantUser: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mockRepo)

			user, err := service.GetUserByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser, user)
			}
		})
	}
}
