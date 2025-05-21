package product

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/grpc_on_http/internal/product/repository"
)

type GrpcController struct {
	UnimplementedProductServiceServer
	productRepository repository.ProductRepository
}

func NewProductGRPCController(repo repository.ProductRepository) ProductServiceServer {
	return &GrpcController{productRepository: repo}
}

func (s *GrpcController) GetByID(ctx context.Context, req *GetByIDRequest) (*GetByIDResponse, error) {
	p, err := s.productRepository.GetProductByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &GetByIDResponse{
		Name:  p.Name,
		Price: p.Price,
		Image: p.Image,
	}, nil
}
