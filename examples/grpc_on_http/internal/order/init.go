package order

import (
	"github.com/ciazhar/go-start-small/examples/grpc_on_http/internal/order/controller"
	"github.com/ciazhar/go-start-small/examples/grpc_on_http/internal/product/controller/grpc/product"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func Init(fiber *fiber.App, grpcClient *grpc.ClientConn) {
	productClient := product.NewProductServiceClient(grpcClient)
	orderController := controller.NewOrderController(productClient)

	fiber.Get("/order", orderController.GetOrderDetail)
}
