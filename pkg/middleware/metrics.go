package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/kxddry/avito-backend-internship-2025/pkg/metrics"
)

func Metrics() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			path := c.Path()

			metrics.HTTPRequestsInFlight.WithLabelValues(path).Inc()
			defer metrics.HTTPRequestsInFlight.WithLabelValues(path).Dec()

			err := next(c)

			duration := time.Since(start).Seconds()
			status := c.Response().Status
			method := c.Request().Method

			metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
			metrics.HTTPRequestsTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()

			return err
		}
	}
}
