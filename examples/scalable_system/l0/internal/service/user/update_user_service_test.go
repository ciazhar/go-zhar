package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	mock_user "github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/repository/postgres/user/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRepositoryContract(ctrl)
	service := NewUserService(mockRepo)

	mockID := uuid.New().String()
	mockReq := request.UpdateUserBodyRequest{
		FullName: "Updated Name",
		Email:    "updated@example.com",
	}

	tests := []struct {
		name     string
		id       string
		input    request.UpdateUserBodyRequest
		mockFunc func(mockRepo *mock_user.MockUserRepositoryContract)
		wantErr  bool
	}{
		{
			name:  "success update user",
			id:    mockID,
			input: mockReq,
			mockFunc: func(mockRepo *mock_user.MockUserRepositoryContract) {
				mockRepo.
					EXPECT().
					UpdateUser(gomock.Any(), mockID, mockReq).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "failed update user",
			id:    mockID,
			input: mockReq,
			mockFunc: func(mockRepo *mock_user.MockUserRepositoryContract) {
				mockRepo.
					EXPECT().
					UpdateUser(gomock.Any(), mockID, mockReq).
					Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mockRepo)

			err := service.UpdateUser(context.Background(), tt.id, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
