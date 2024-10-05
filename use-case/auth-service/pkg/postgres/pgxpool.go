package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

func InitPostgres(ctx context.Context) *pgxpool.Pool {
	var err error
	dbPool, err := pgxpool.Connect(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	return dbPool
}
