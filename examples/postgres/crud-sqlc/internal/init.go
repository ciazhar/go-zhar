package internal

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/postgres/crud-sqlc/internal/controller"
	db "github.com/ciazhar/go-zhar/examples/postgres/crud-sqlc/internal/generated/repository"
	"github.com/ciazhar/go-zhar/examples/postgres/crud-sqlc/internal/service"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Init(ctx context.Context, app *fiber.App, queries *db.Queries, db *pgxpool.Pool, logger *logger.Logger) {

	s := service.NewProductService(queries, db, logger)
	c := controller.NewProductController(ctx, s)

	// Setup routes
	api := app.Group("/api")

	// Products routes
	products := api.Group("/products")
	products.Post("/", c.CreateProduct)
	products.Get("/", c.GetProducts)
	products.Get("/cursor", c.GetProductsCursor)
	products.Put("/:id/price", c.UpdateProductPrice)
	products.Delete("/:id", c.DeleteProduct)

}
