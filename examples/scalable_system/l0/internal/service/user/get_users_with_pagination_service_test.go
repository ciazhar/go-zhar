package user

import (
	"context"
	"errors"
	"testing"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/response"
	mock_user "github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/repository/postgres/user/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRepositoryContract(ctrl)
	service := NewUserService(mockRepo)

	mockReq := request.GetUsersQueryParam{
		Page: 1,
		Size: 10,
	}

	mockUsers := []response.User{
		{
			ID:       "user-1",
			Username: "hafidz",
			Email:    "hafidz@example.com",
			FullName: "Muhammad Hafidz",
		},
		{
			ID:       "user-2",
			Username: "johndoe",
			Email:    "johndoe@example.com",
			FullName: "John Doe",
		},
	}
	mockTotal := int64(2)

	tests := []struct {
		name      string
		input     request.GetUsersQueryParam
		mockFunc  func(mockRepo *mock_user.MockUserRepositoryContract)
		wantUsers []response.User
		wantTotal int64
		wantErr   bool
	}{
		{
			name:  "success fetch users",
			input: mockReq,
			mockFunc: func(mockRepo *mock_user.MockUserRepositoryContract) {
				mockRepo.
					EXPECT().
					GetUsersWithPagination(gomock.Any(), mockReq.Page, mockReq.Size).
					Return(mockUsers, mockTotal, nil)
			},
			wantUsers: mockUsers,
			wantTotal: mockTotal,
			wantErr:   false,
		},
		{
			name:  "failed fetch users",
			input: mockReq,
			mockFunc: func(mockRepo *mock_user.MockUserRepositoryContract) {
				mockRepo.
					EXPECT().
					GetUsersWithPagination(gomock.Any(), mockReq.Page, mockReq.Size).
					Return(nil, int64(0), errors.New("db error"))
			},
			wantUsers: nil,
			wantTotal: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mockRepo)

			users, total, err := service.GetUsersWithPagination(context.Background(), tt.input.Page, tt.input.Size)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, users)
				assert.Equal(t, int64(0), total)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUsers, users)
				assert.Equal(t, tt.wantTotal, total)
			}
		})
	}
}
