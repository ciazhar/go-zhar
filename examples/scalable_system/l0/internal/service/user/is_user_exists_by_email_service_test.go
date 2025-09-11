package user

import (
	"context"
	"errors"
	"testing"

	mock_user "github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/repository/postgres/user/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestIsUserExistsByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRepositoryContract(ctrl)
	service := NewUserService(mockRepo)

	mockEmail := "hafidz@example.com"

	tests := []struct {
		name      string
		email     string
		mockFunc  func(mockRepo *mock_user.MockUserRepositoryContract)
		wantExist bool
		wantErr   bool
	}{
		{
			name:  "user exists",
			email: mockEmail,
			mockFunc: func(mockRepo *mock_user.MockUserRepositoryContract) {
				mockRepo.
					EXPECT().
					IsUserExistsByEmail(gomock.Any(), mockEmail).
					Return(true, nil)
			},
			wantExist: true,
			wantErr:   false,
		},
		{
			name:  "user does not exist",
			email: mockEmail,
			mockFunc: func(mockRepo *mock_user.MockUserRepositoryContract) {
				mockRepo.
					EXPECT().
					IsUserExistsByEmail(gomock.Any(), mockEmail).
					Return(false, nil)
			},
			wantExist: false,
			wantErr:   false,
		},
		{
			name:  "db error when checking user",
			email: mockEmail,
			mockFunc: func(mockRepo *mock_user.MockUserRepositoryContract) {
				mockRepo.
					EXPECT().
					IsUserExistsByEmail(gomock.Any(), mockEmail).
					Return(false, errors.New("db error"))
			},
			wantExist: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mockRepo)

			exists, err := service.IsUserExistsByEmail(context.Background(), tt.email)

			if tt.wantErr {
				assert.Error(t, err)
				assert.False(t, exists)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantExist, exists)
			}
		})
	}
}
