package user

import (
	"context"
	"errors"
	"testing"

	"github.com/ciazhar/go-zhar-scalable-system-l0/internal/model/request"
	mockpostgres "github.com/ciazhar/go-zhar/pkg/postgres/mock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpsertUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mockpostgres.NewMockPgxPool(ctrl)
	repo := NewUserRepository(mock)
	mockReq := request.UpsertUserBodyRequest{
		Id:       "11111111-1111-1111-1111-111111111111",
		Username: "upsertuser",
		Email:    "upsert@example.com",
		Password: "secret",
		FullName: "Upsert User",
	}

	tests := []struct {
		name     string
		req      request.UpsertUserBodyRequest
		mockFunc func(mockPool *mockpostgres.MockPgxPool)
		wantErr  bool
	}{
		{
			name: "success - user upserted",
			req:  mockReq,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				mockPool.
					EXPECT().
					Exec(gomock.Any(), queryUpsertUser, mockReq.Id, mockReq.Username, mockReq.Email, mockReq.Password, mockReq.FullName).
					Return(pgconn.NewCommandTag("INSERT 1"), nil)
			},
			wantErr: false,
		},
		{
			name: "failed - db error",
			req:  mockReq,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				mockPool.
					EXPECT().
					Exec(gomock.Any(), queryUpsertUser, mockReq.Id, mockReq.Username, mockReq.Email, mockReq.Password, mockReq.FullName).
					Return(pgconn.CommandTag{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mock)

			err := repo.UpsertUserByID(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
