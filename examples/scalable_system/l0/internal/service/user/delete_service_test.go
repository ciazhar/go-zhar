package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"

	mock_user "github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/repository/postgres/user/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRepositoryContract(ctrl)
	service := NewUserService(mockRepo)
	mockID := uuid.New().String()

	tests := []struct {
		name     string
		id       string
		mockFunc func(mockRepo *mock_user.MockUserRepositoryContract)
		wantErr  bool
	}{
		{
			name: "success delete user",
			id:   mockID,
			mockFunc: func(mockRepo *mock_user.MockUserRepositoryContract) {
				mockRepo.
					EXPECT().
					DeleteUser(gomock.Any(), mockID).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "failed delete user",
			id:   mockID,
			mockFunc: func(mockRepo *mock_user.MockUserRepositoryContract) {
				mockRepo.
					EXPECT().
					DeleteUser(gomock.Any(), mockID).
					Return(errors.New("db delete error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mockRepo)

			err := service.DeleteUser(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
