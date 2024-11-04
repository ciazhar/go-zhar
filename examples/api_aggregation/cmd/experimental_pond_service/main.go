package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alitto/pond/v2"
	"github.com/ciazhar/go-start-small/examples/api_aggregation/pkg"
	"github.com/gofiber/fiber/v2"
)

// RequestConfig holds the configuration for an API request
type RequestConfig struct {
	URL         string
	Target      interface{}
	ServiceName string
}

// makeRequest creates a function that performs an HTTP request and decodes the response
func makeRequest(cfg RequestConfig) error {

	req, err := http.NewRequest(http.MethodGet, cfg.URL, nil)
	if err != nil {
		return fmt.Errorf("%s: failed to create request: %w", cfg.ServiceName, err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s: request failed: %w", cfg.ServiceName, err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: failed to read response body: %w", cfg.ServiceName, err)
	}

	// fmt.Printf("Raw response from %s: %s\n", cfg.ServiceName, string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: unexpected status code %d: %s",
			cfg.ServiceName, resp.StatusCode, string(bodyBytes))
	}

	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(cfg.Target); err != nil {
		return fmt.Errorf("%s: failed to decode response: %w\nResponse body: %s",
			cfg.ServiceName, err, string(bodyBytes))
	}

	// fmt.Printf("Decoded %s response: %+v\n", cfg.ServiceName, cfg.Target)

	return nil
}

func main() {
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	app.Get("/dashboard", getDashboard)

	if err := app.Listen(":3004"); err != nil {
		panic(fmt.Sprintf("failed to start server: %v", err))
	}
}

func getDashboard(c *fiber.Ctx) error {
	userID := c.Query("userID")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userID is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 4*time.Second)
	defer cancel()

	data, err := getDashboardData(ctx, userID)
	if err != nil {
		return errorHandler(c, err)
	}

	return c.JSON(data)
}

func getDashboardData(ctx context.Context, userID string) (*pkg.GetDashboardDataResponse, error) {
	dashboard := &pkg.GetDashboardDataResponse{
		Orders:   make([]pkg.GetUserOrdersResponse, 0),
		Products: make([]pkg.GetProductRecommendationsResponse, 0),
	}

	requests := []RequestConfig{
		{
			URL:         fmt.Sprintf("http://localhost:3001/orders?userID=%s", userID),
			Target:      &dashboard.Orders,
			ServiceName: "OrderService",
		},
		{
			URL:         fmt.Sprintf("http://localhost:3002/recommendations?userID=%s", userID),
			Target:      &dashboard.Products,
			ServiceName: "ProductService",
		},
		{
			URL:         fmt.Sprintf("http://localhost:3003/profile?userID=%s", userID),
			Target:      &dashboard.User,
			ServiceName: "UserService",
		},
	}

	// Set up pond worker pool
	pool := pond.NewPool(len(requests)) // Adjust worker and queue size as needed
	defer pool.StopAndWait()

	// errorChan := make(chan error, len(requests))
	var errors []error

	// Submit tasks to the pond pool
	for _, req := range requests {
		// req := req
		task := pool.SubmitErr(func() error {
			return makeRequest(req)
		})
		err := task.Wait()
		if err != nil {
			errors = append(errors, err)
		}
	}

	// Wait for all tasks to complete
	pool.StopAndWait()

	if len(errors) > 0 {
		return nil, fmt.Errorf("failed to fetch dashboard data: %v", errors)
	}

	return dashboard, nil
}

// Middleware for common error handling
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if err == context.DeadlineExceeded {
		code = fiber.StatusGatewayTimeout
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
