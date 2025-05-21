package product

import (
	"github.com/ciazhar/go-start-small/examples/grpc_on_http/internal/product/controller/grpc/product"
	"github.com/ciazhar/go-start-small/examples/grpc_on_http/internal/product/repository"
	"google.golang.org/grpc"
)

func Init(grpcServer *grpc.Server) {
	productRepo := repository.NewDummyProductRepository()
	controller := product.NewProductGRPCController(productRepo)
	product.RegisterProductServiceServer(grpcServer, controller)
}
