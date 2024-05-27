package order

import (
	"github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/internal/order/controller"
	"github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/internal/order/repository"
	"github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/internal/order/service"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func Init(app *fiber.App, httpClient *http.Client, tracer trace.Tracer) {

	r := repository.NewOrderRepository(tracer)
	hr := repository.NewUserHTTPRepository(httpClient, tracer)
	s := service.NewOrderService(r, hr, tracer)
	c := controller.NewOrderController(s, tracer)

	app.Post("/orders", c.AddOrder)
	app.Get("/orders/:order_id", c.GetOrderByOrderID)
	app.Get("/orders", c.GetAllOrders)
	app.Delete("/orders/:order_id", c.DeleteOrder)
	app.Put("/orders", c.UpdateOrder)
}
