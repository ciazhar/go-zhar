package service

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/mock_testing/internal/model"
	"github.com/ciazhar/go-start-small/examples/mock_testing/internal/repository"
)

type AccidentReportService struct {
	repo repository.AccidentReportRepositoryInterface
}

type AccidentReportServiceInterface interface {
	GetAccidentReport(ctx context.Context, id string) (*model.AccidentReport, error)
}

func NewAccidentReportService(repo repository.AccidentReportRepositoryInterface) *AccidentReportService {
	return &AccidentReportService{
		repo: repo,
	}
}

func (s *AccidentReportService) GetAccidentReport(ctx context.Context, id string) (*model.AccidentReport, error) {
	return s.repo.GetAccidentReport(ctx, id)
}
