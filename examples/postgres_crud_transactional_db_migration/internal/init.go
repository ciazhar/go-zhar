package internal

import (
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/controller"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/repository"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(app *fiber.App, pool *pgxpool.Pool) {
	ir := repository.NewPgxInventoryRepository(pool)
	or := repository.NewPgxOrderRepository(pool)
	pr := repository.NewPgxPaymentRepository(pool)
	os := service.NewOrderService(or, ir, pr)

	oc := controller.NewOrderController(os)

	app.Get("/orders", oc.GetAllOrders)

	app.Get("/orders/:id", oc.GetOrder)

	app.Delete("/orders/:id", oc.Delete)

	app.Post("/orders", oc.CreateOrder)
}
