package order

import (
	"net/http"

	"github.com/ciazhar/go-start-small/examples/http_distributed_tracing/internal/order/controller"
	"github.com/ciazhar/go-start-small/examples/http_distributed_tracing/internal/order/repository"
	"github.com/ciazhar/go-start-small/examples/http_distributed_tracing/internal/order/service"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App, httpClient *http.Client) {

	r := repository.NewOrderRepository()
	hr := repository.NewUserHTTPRepository(httpClient)
	s := service.NewOrderService(r, hr)
	c := controller.NewOrderController(s)

	app.Post("/orders", c.AddOrder)
	app.Get("/orders/:order_id", c.GetOrderByOrderID)
	app.Get("/orders", c.GetAllOrders)
	app.Delete("/orders/:order_id", c.DeleteOrder)
	app.Put("/orders", c.UpdateOrder)
}
