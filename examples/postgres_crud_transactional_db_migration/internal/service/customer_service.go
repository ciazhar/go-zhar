package service

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/repository"
)

type CustomerService struct {
	Repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) *CustomerService {
	return &CustomerService{Repo: repo}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, name, email string) (int, error) {
	return s.Repo.CreateCustomer(ctx, name, email)
}
