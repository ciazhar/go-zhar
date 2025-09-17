package user

import (
	"context"
	"errors"
	"testing"

	mockpostgres "github.com/ciazhar/go-zhar/pkg/postgres/mock"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mockpostgres.NewMockPgxPool(ctrl)
	repo := NewUserRepository(mock)
	mockReq := request.CreateUserBodyRequest{
		Username: "hafidz",
		Email:    "hafidz@example.com",
		Password: "hashed-password",
		FullName: "Muhammad Hafidz",
	}

	tests := []struct {
		name     string
		input    request.CreateUserBodyRequest
		mockFunc func(mockPool *mockpostgres.MockPgxPool)
		wantErr  bool
	}{
		{
			name:  "success insert user",
			input: mockReq,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				mockPool.
					EXPECT().
					Exec(gomock.Any(), queryCreateUser, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(pgconn.CommandTag{}, nil)
			},
			wantErr: false,
		},
		{
			name:  "failed insert user",
			input: mockReq,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				mockPool.
					EXPECT().
					Exec(gomock.Any(), queryCreateUser, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(pgconn.CommandTag{}, errors.New("db insert error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mock)

			err := repo.CreateUser(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
