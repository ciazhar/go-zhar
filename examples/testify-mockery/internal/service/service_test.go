package service_test

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/testify-mockery/internal/mocks"
	"github.com/ciazhar/go-zhar/examples/testify-mockery/internal/model"
	"github.com/ciazhar/go-zhar/examples/testify-mockery/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetAccidentReport(t *testing.T) {
	mockRepo := mocks.NewRepositoryInterface(t)
	svc := service.NewService(mockRepo)

	reportID := "1234"
	expectedReport := &model.AccidentReport{
		ID:     reportID,
		Report: "Accident details",
	}

	mockRepo.On("GetAccidentReport", mock.Anything, reportID).Return(expectedReport, nil)

	ctx := context.Background()
	report, err := svc.GetAccidentReport(ctx, reportID)

	assert.NoError(t, err)
	assert.Equal(t, expectedReport, report)

	mockRepo.AssertExpectations(t)
}
