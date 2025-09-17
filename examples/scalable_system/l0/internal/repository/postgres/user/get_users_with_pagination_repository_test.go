package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/response"
	"github.com/ciazhar/go-zhar/pkg/postgres"
	mockpostgres "github.com/ciazhar/go-zhar/pkg/postgres/mock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mockpostgres.NewMockPgxPool(ctrl)
	repo := NewUserRepository(mock)
	mockTime := time.Now()
	mockRes := []response.User{
		{ID: "1", Username: "alice", Email: "alice@mail.com", FullName: "Alice", CreatedAt: mockTime, UpdatedAt: mockTime},
		{ID: "2", Username: "bob", Email: "bob@mail.com", FullName: "Bob", CreatedAt: mockTime, UpdatedAt: mockTime},
	}
	mockRows := postgres.MockRows{
		Rows: []postgres.MockRow{
			{Values: []any{mockRes[0].ID, mockRes[0].Username, mockRes[0].Email, mockRes[0].FullName, mockRes[0].CreatedAt, mockRes[0].UpdatedAt}},
			{Values: []any{mockRes[1].ID, mockRes[1].Username, mockRes[1].Email, mockRes[1].FullName, mockRes[1].CreatedAt, mockRes[1].UpdatedAt}},
		},
	}
	mockRows2 := postgres.MockRows{
		Rows: []postgres.MockRow{
			{Values: []any{mockRes[0].ID, mockRes[0].Username, mockRes[0].Email, mockRes[0].FullName, mockRes[0].CreatedAt, mockRes[0].UpdatedAt}},
		},
	}
	errMockRows := postgres.MockRows{
		Rows: []postgres.MockRow{
			{Values: []any{"invalid-id-only"}}, // not enough fields → Scan error
		},
	}

	tests := []struct {
		name      string
		page      int
		limit     int
		mockFunc  func(mockPool *mockpostgres.MockPgxPool)
		want      []response.User
		wantTotal int64
		wantErr   bool
	}{
		{
			name:  "success - multiple users",
			page:  1,
			limit: 10,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				// Mock Query for users
				mockPool.
					EXPECT().
					Query(gomock.Any(), queryGetUsersWithPagination, 10, 0).
					Return(&mockRows, nil)

				// Mock QueryRow for total count
				mockPool.
					EXPECT().
					QueryRow(gomock.Any(), queryCountUsers).
					Return(pgx.Row(postgres.MockRow{Values: []any{int64(2)}}))
			},
			want:      mockRes,
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name:  "failed - query error",
			page:  1,
			limit: 10,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				mockPool.
					EXPECT().
					Query(gomock.Any(), queryGetUsersWithPagination, 10, 0).
					Return(nil, errors.New("db error"))
			},
			want:      nil,
			wantTotal: 0,
			wantErr:   true,
		},
		{
			name:  "failed - scan error",
			page:  1,
			limit: 10,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				// Mock Query returning row with mismatched values
				mockPool.
					EXPECT().
					Query(gomock.Any(), queryGetUsersWithPagination, 10, 0).
					Return(&errMockRows, nil)

				// QueryRow should not be called in this scenario, so no expectation
			},
			want:      nil,
			wantTotal: 0,
			wantErr:   true,
		},
		{
			name:  "failed - count error",
			page:  1,
			limit: 10,
			mockFunc: func(mockPool *mockpostgres.MockPgxPool) {
				// Return valid rows first
				mockPool.
					EXPECT().
					Query(gomock.Any(), queryGetUsersWithPagination, 10, 0).
					Return(&mockRows2, nil)

				// Fail on QueryRow (count users)
				mockPool.
					EXPECT().
					QueryRow(gomock.Any(), queryCountUsers).
					Return(pgx.Row(postgres.MockRow{Err: errors.New("count error")}))
			},
			want:      nil,
			wantTotal: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mock)

			got, total, err := repo.GetUsersWithPagination(context.Background(), tt.page, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.Equal(t, int64(0), total)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
				assert.Equal(t, tt.wantTotal, total)
			}
		})
	}
}
