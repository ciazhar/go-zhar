package middleware

import (
	metrics "github.com/ciazhar/go-zhar/pkg/prometheus"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// lanjut ke handler berikut
		err := c.Next()

		// label
		status := strconv.Itoa(c.Response().StatusCode())
		method := c.Method()
		path := c.Route().Path

		// update metrics
		metrics.HttpRequestsTotal.WithLabelValues(method, path, status).Inc()
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(method, path).Observe(duration)

		return err
	}
}
