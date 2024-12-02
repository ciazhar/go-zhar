package repository

import (
	"context"

	"github.com/ciazhar/go-start-small/examples/jwt_authorization_service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersPostgresRepositoryInterface interface {
	GetByUsername(ctx context.Context, username string) (model.User, error)
	Insert(ctx context.Context, tx pgx.Tx, user model.User) error
	BeginTransaction(ctx context.Context) (pgx.Tx, error)
}

type UsersPostgresRepository struct {
	pg *pgxpool.Pool
}

func NewUsersPostgresRepository(pg *pgxpool.Pool) *UsersPostgresRepository {
	return &UsersPostgresRepository{pg: pg}
}

func (r *UsersPostgresRepository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User
	err := r.pg.QueryRow(ctx, "SELECT id, username, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UsersPostgresRepository) Insert(ctx context.Context, tx pgx.Tx, user model.User) error {
	_, err := tx.Exec(ctx, "INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersPostgresRepository) BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	return r.pg.BeginTx(ctx, pgx.TxOptions{})
}
