package main

import (
	"github.com/ciazhar/go-start-small/examples/grpc_on_http/internal/order"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {

	app := fiber.New()

	//for insecure
	grpcConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer grpcConn.Close()

	// for tls
	//creds := credentials.NewClientTLSFromCert(nil, "")
	//grpcConn, err := grpc.NewClient("host:port", grpc.WithTransportCredentials(creds))
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//defer grpcConn.Close()

	order.Init(app, grpcConn)

	log.Println("server running on port 3000")
	err = app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}

}
