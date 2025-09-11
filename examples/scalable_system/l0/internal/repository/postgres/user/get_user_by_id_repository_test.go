package user

import (
	"context"
	"errors"
	"github.com/ciazhar/go-zhar/pkg/postgres"
	"github.com/google/uuid"
	"testing"
	"time"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/response"
	mockpostgres "github.com/ciazhar/go-zhar/pkg/postgres/mock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mockpostgres.NewMockPgxPool(ctrl)
	repo := NewUserRepository(mock)
	mockReq := uuid.New().String()
	mockTime := time.Now()
	mockRes := response.User{
		ID:        mockReq,
		Username:  "testuser",
		Email:     "test@example.com",
		FullName:  "Test User",
		CreatedAt: mockTime,
		UpdatedAt: mockTime,
	}
	mockRow := pgx.Row(postgres.MockRow{
		Values: []interface{}{
			mockRes.ID,
			mockRes.Username,
			mockRes.Email,
			mockRes.FullName,
			mockRes.CreatedAt,
			mockRes.UpdatedAt,
		},
		Err: nil,
	})

	tests := []struct {
		name     string
		id       string
		mockFunc func(mockPool *mockpostgres.MockPgxPool)
		want     *response.User
		wantErr  bool
	}{
		{
			name: "success - user found",
			id:   mockReq,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				mockPool.
					EXPECT().
					QueryRow(gomock.Any(), queryGetUserByID, mockReq).
					Return(mockRow)
			},
			want:    &mockRes,
			wantErr: false,
		},
		{
			name: "failed - db error",
			id:   mockReq,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				mockPool.
					EXPECT().
					QueryRow(gomock.Any(), queryGetUserByID, mockReq).
					Return(pgx.Row(postgres.MockRow{Err: errors.New("db error")}))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mock)

			got, err := repo.GetUserByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
