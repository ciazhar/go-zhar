package repository

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/testify-mockery/internal/model"
)

type Repository struct {
	a string
}

type RepositoryInterface interface {
	GetAccidentReport(ctx context.Context, id string) (*model.AccidentReport, error)
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetAccidentReport(ctx context.Context, id string) (*model.AccidentReport, error) {
	return &model.AccidentReport{
		ID:     id,
		Report: "report",
	}, nil
}
