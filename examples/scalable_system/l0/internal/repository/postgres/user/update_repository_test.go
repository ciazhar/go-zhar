package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	mockpostgres "github.com/ciazhar/go-zhar/pkg/postgres/mock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mockpostgres.NewMockPgxPool(ctrl)
	repo := NewUserRepository(mock)
	mockId := uuid.New().String()
	mockReq := request.UpdateUserBodyRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
		FullName: "Updated User",
	}

	tests := []struct {
		name     string
		id       string
		req      request.UpdateUserBodyRequest
		mockFunc func(mockPool *mockpostgres.MockPgxPool)
		wantErr  bool
	}{
		{
			name: "success - user updated",
			id:   mockId,
			req:  mockReq,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				mockPool.
					EXPECT().
					Exec(gomock.Any(), queryUpdateUser, mockReq.Username, mockReq.Email, mockReq.FullName, mockId).
					Return(pgconn.NewCommandTag("UPDATE 1"), nil)
			},
			wantErr: false,
		},
		{
			name: "failed - db error",
			id:   mockId,
			req:  mockReq,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				mockPool.
					EXPECT().
					Exec(gomock.Any(), queryUpdateUser, mockReq.Username, mockReq.Email, mockReq.FullName, mockId).
					Return(pgconn.CommandTag{}, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "failed - no rows updated",
			id:   mockId,
			req:  mockReq,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				mockPool.
					EXPECT().
					Exec(gomock.Any(), queryUpdateUser, mockReq.Username, mockReq.Email, mockReq.FullName, mockId).
					Return(pgconn.NewCommandTag("UPDATE 0"), nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mock)

			err := repo.UpdateUser(context.Background(), tt.id, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
