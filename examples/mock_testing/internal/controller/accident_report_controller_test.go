package controller

import (
	"errors"
	"github.com/ciazhar/go-start-small/examples/mock_testing/internal/mocks"
	"github.com/ciazhar/go-start-small/examples/mock_testing/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetAccidentReportHandler tests the GetAccidentReportHandler
func TestGetAccidentReportHandler(t *testing.T) {
	s := mocks.NewAccidentReportServiceInterface(t)
	controller := NewAccidentReportController(s)

	app := fiber.New()
	app.Get("/accident-report/:id", controller.GetAccidentReportHandler)

	tests := []struct {
		description  string
		route        string
		mockSetup    func()
		expectedCode int
		expectedBody string
	}{
		{
			description: "successful report retrieval",
			route:       "/accident-report/123",
			mockSetup: func() {
				s.On("GetAccidentReport", mock.Anything, "123").Return(&model.AccidentReport{
					ID: "123", Report: "Report Details"}, nil).Once()
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"123","report":"Report Details"}`,
		},
		{
			description: "internal server error",
			route:       "/accident-report/456",
			mockSetup: func() {
				s.On("GetAccidentReport", mock.Anything, "456").Return(nil, errors.New("some error")).Once()
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"Failed to get accident report"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tt.mockSetup()
			req := httptest.NewRequest(http.MethodGet, tt.route, nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedCode, resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)
			assert.JSONEq(t, tt.expectedBody, string(body))
		})
	}
}
