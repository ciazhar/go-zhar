package user

import (
	"context"
	"errors"
	"github.com/ciazhar/go-zhar/pkg/postgres"
	mockpostgres "github.com/ciazhar/go-zhar/pkg/postgres/mock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestExistsByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mockpostgres.NewMockPgxPool(ctrl)
	repo := NewUserRepository(mock)
	mockReq := "test@example.com"

	tests := []struct {
		name     string
		email    string
		mockFunc func(mockPool *mockpostgres.MockPgxPool)
		want     bool
		wantErr  bool
	}{
		{
			name:  "success - user exists",
			email: mockReq,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				mockPool.
					EXPECT().
					QueryRow(gomock.Any(), queryIsUserExistsByEmail, mockReq).
					Return(pgx.Row(postgres.MockRow{Values: []any{true}, Err: nil}))
			},
			want:    true,
			wantErr: false,
		},
		{
			name:  "success - user does not exist",
			email: mockReq,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				mockPool.
					EXPECT().
					QueryRow(gomock.Any(), queryIsUserExistsByEmail, mockReq).
					Return(pgx.Row(postgres.MockRow{Values: []any{
						false,
					}, Err: nil,
					}))
			},
			want:    false,
			wantErr: false,
		},
		{
			name:  "db error",
			email: mockReq,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				mockPool.
					EXPECT().
					QueryRow(gomock.Any(), queryIsUserExistsByEmail, mockReq).
					Return(pgx.Row(postgres.MockRow{Err: errors.New("db error")}))
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mock)

			got, err := repo.IsUserExistsByEmail(context.Background(), tt.email)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
