package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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
func makeRequest(cfg RequestConfig) pkg.GenericFunction {
    return func(ctx context.Context) error {
        req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.URL, nil)
        if err != nil {
            return fmt.Errorf("%s: failed to create request: %w", cfg.ServiceName, err)
        }

        req.Header.Set("Accept", "application/json")
        
        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            return fmt.Errorf("%s: request failed: %w", cfg.ServiceName, err)
        }
        defer resp.Body.Close()

        // Debug: Print raw response
        bodyBytes, err := io.ReadAll(resp.Body)
        if err != nil {
            return fmt.Errorf("%s: failed to read response body: %w", cfg.ServiceName, err)
        }

        fmt.Printf("Raw response from %s: %s\n", cfg.ServiceName, string(bodyBytes))

        if resp.StatusCode != http.StatusOK {
            return fmt.Errorf("%s: unexpected status code %d: %s",
                cfg.ServiceName, resp.StatusCode, string(bodyBytes))
        }

        // Create new reader from bytes since we consumed the original body
        if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(cfg.Target); err != nil {
            return fmt.Errorf("%s: failed to decode response: %w\nResponse body: %s", 
                cfg.ServiceName, err, string(bodyBytes))
        }

        // Debug: Print decoded struct
        fmt.Printf("Decoded %s response: %+v\n", cfg.ServiceName, cfg.Target)

        return nil
    }
}

func main() {
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	app.Get("/dashboard", getDashboard)

	if err := app.Listen(":3000"); err != nil {
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

	// Create context with timeout
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

	// Configure requests
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

	// Create functions for each request
	funcs := make([]pkg.GenericFunction, len(requests))
	for i, req := range requests {
		funcs[i] = makeRequest(req)
	}

	// Configure async runner options
	opts := pkg.Options{
		Timeout:         3 * time.Second,
		MaxConcurrent:   3,
		ContinueOnError: false, // Stop on first error
	}

	// Execute requests concurrently
	results := pkg.RunAsync(ctx, opts, funcs...)

	// Check for errors
	var errors []error
	for _, result := range results {
		if result.Err != nil {
			errors = append(errors, result.Err)
		}
	}

	// If any requests failed, return all errors
	if len(errors) > 0 {
		return nil, fmt.Errorf("failed to fetch dashboard data: %v", errors)
	}

	return dashboard, nil
}

// Middleware for common error handling
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	// Check for specific error types
	if err == context.DeadlineExceeded {
		code = fiber.StatusGatewayTimeout
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
