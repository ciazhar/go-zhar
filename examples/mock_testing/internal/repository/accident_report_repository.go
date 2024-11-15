package repository

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/mock_testing/internal/model"
)

type AccidentReportRepository struct {
}

type AccidentReportRepositoryInterface interface {
	GetAccidentReport(ctx context.Context, id string) (*model.AccidentReport, error)
}

func NewAccidentReportRepository() *AccidentReportRepository {
	return &AccidentReportRepository{}
}

func (r *AccidentReportRepository) GetAccidentReport(ctx context.Context, id string) (*model.AccidentReport, error) {
	return &model.AccidentReport{
		ID:     id,
		Report: "report",
	}, nil
}
