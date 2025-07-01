package main

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-start-small/examples/db_pg_csv_zip_http/internal"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	fmt.Println("starting data generator")

	// Connect to PostgreSQL
	dbUrl := "postgres://user:password@localhost:5432/demo"
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	fmt.Println("data generator started")

	internal.InitDataGenerator(pool)

	fmt.Println("data generator finished")
}
