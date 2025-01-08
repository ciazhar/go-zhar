package internal

import (
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/controller"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/repository"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(app *fiber.App, pool *pgxpool.Pool) {
	customerRepository := repository.NewCustomerRepository(pool)
	productRepository := repository.NewProductRepository(pool)
	orderRepository := repository.NewOrderRepository(pool)
	paymentRepository := repository.NewPaymentRepository(pool)
	shipmentRepository := repository.NewShipmentRepository(pool)

	customerService := service.NewCustomerService(customerRepository)
	productService := service.NewProductService(productRepository)
	orderService := service.NewOrderService(orderRepository, productRepository, paymentRepository, shipmentRepository)

	customerController := controller.NewCustomerController(customerService)
	productController := controller.NewProductController(productService)
	orderController := controller.NewOrderController(orderService)

	customers := app.Group("/customers")
	customers.Post("/", customerController.CreateCustomer)

	products := app.Group("/products")
	products.Get("/", productController.GetProducts)
	products.Post("/", productController.CreateProduct)

	orders := app.Group("/orders")
	orders.Post("/", orderController.PlaceOrder)
	orders.Post("/payment", orderController.ProcessPayment)
	orders.Post("/ship", orderController.ShipOrder)
	orders.Patch("/:orderID/delivered", orderController.MarkOrderDelivered)

}
