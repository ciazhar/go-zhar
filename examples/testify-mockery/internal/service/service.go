package service

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/testify-mockery/internal/model"
	"github.com/ciazhar/go-zhar/examples/testify-mockery/internal/repository"
)

type Service struct {
	repo repository.RepositoryInterface
}

func NewService(repo repository.RepositoryInterface) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAccidentReport(ctx context.Context, id string) (*model.AccidentReport, error) {
	return s.repo.GetAccidentReport(ctx, id)
}
