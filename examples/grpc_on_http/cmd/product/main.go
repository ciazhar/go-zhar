package main

import (
	"github.com/ciazhar/go-start-small/examples/grpc_on_http/internal/product"
	grpc_middleware "github.com/ciazhar/go-start-small/examples/grpc_on_http/pkg/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	apiKey := "test-api-key"
	gm := grpc_middleware.New(apiKey)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(gm.UnaryInterceptor))
	reflection.Register(grpcServer)

	product.Init(grpcServer)

	log.Println("gRPC server running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}
