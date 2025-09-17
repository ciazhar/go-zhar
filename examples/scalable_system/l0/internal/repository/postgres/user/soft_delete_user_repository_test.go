package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"

	mockpostgres "github.com/ciazhar/go-zhar/pkg/postgres/mock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mockpostgres.NewMockPgxPool(ctrl)
	repo := NewUserRepository(mock)
	mockReq := uuid.New().String()

	tests := []struct {
		name     string
		id       string
		mockFunc func(mockPool *mockpostgres.MockPgxPool)
		wantErr  bool
	}{
		{
			name: "success delete user",
			id:   mockReq,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				// Simulate 1 row affected
				mockPool.
					EXPECT().
					Exec(gomock.Any(), querySoftDeleteUser, gomock.Any()).
					Return(pgconn.NewCommandTag("UPDATE 1"), nil)
			},
			wantErr: false,
		},
		{
			name: "failed delete user - db error",
			id:   mockReq,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				mockPool.
					EXPECT().
					Exec(gomock.Any(), querySoftDeleteUser, gomock.Any()).
					Return(pgconn.CommandTag{}, errors.New("db delete error"))
			},
			wantErr: true,
		},
		{
			name: "failed delete user - no rows deleted",
			id:   mockReq,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				// Simulate 0 rows affected
				mockPool.
					EXPECT().
					Exec(gomock.Any(), querySoftDeleteUser, gomock.Any()).
					Return(pgconn.NewCommandTag("UPDATE 0"), nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mock)

			err := repo.SoftDeleteUser(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
