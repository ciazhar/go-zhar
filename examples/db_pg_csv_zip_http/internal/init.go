package internal

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-start-small/examples/db_pg_csv_zip_http/internal/controller"
	"github.com/ciazhar/go-start-small/examples/db_pg_csv_zip_http/internal/repository"
	"github.com/ciazhar/go-start-small/examples/db_pg_csv_zip_http/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitClient(fiber *fiber.App) {
	c := controller.NewEmailController(nil)
	fiber.Post("/email", c.SendEmail)
}

func InitDataGenerator(pool *pgxpool.Pool) {
	r := repository.NewRepository(pool)
	s := service.NewService(r)
	err := s.ImportCSV(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

func InitReport(fiber *fiber.App, pool *pgxpool.Pool) {
	r := repository.NewRepository(pool)
	s := service.NewService(r)
	c := controller.NewEmailController(s)
	fiber.Get("/optimized", c.ExportOptimized)
	fiber.Get("/unoptimized", c.ExportUnoptimized)
}
