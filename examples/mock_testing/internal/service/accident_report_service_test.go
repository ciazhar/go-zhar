package service

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/mock_testing/internal/mocks"
	"github.com/ciazhar/go-start-small/examples/mock_testing/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestService_GetAccidentReport(t *testing.T) {
	mockRepo := mocks.NewAccidentReportRepositoryInterface(t)
	svc := NewAccidentReportService(mockRepo)

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
