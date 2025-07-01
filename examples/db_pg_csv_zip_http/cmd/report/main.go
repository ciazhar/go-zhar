package main

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-start-small/examples/db_pg_csv_zip_http/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "net/http/pprof"
)

func main() {

	fmt.Println("starting report server")

	ctx := context.Background()

	fmt.Println("report server started")

	// Connect to PostgreSQL
	dbUrl := "postgres://user:password@localhost:5432/demo"
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	app := fiber.New()
	internal.InitReport(app, pool)

	err = app.Listen(":3001")
	if err != nil {
		fmt.Println(err)
	}
}
