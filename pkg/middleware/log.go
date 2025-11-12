package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const httpServer = "http_server"

func Zerolog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			req := c.Request()
			res := c.Response()
			start := time.Now()

			reqID := req.Header.Get(echo.HeaderXRequestID)
			if reqID == "" {
				reqID = res.Header().Get(echo.HeaderXRequestID)
			}

			// Выполняем хендлер
			err := next(c)

			latency := time.Since(start)
			status := res.Status

			l := log.Logger.With().
				Str("request_id", reqID).
				Str("remote_ip", c.RealIP()).
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Str("host", req.Host).
				Str("user_agent", req.UserAgent()).
				Int("status", status).
				Dur("latency", latency).
				Int64("bytes_out", res.Size).
				Str("route", req.URL.Path).
				Logger()

			switch {
			case status >= 500:
				if err != nil {
					l.Error().Err(err).Msg(httpServer)
				} else {
					l.Error().Msg(httpServer)
				}
			case status >= 400:
				if err != nil {
					l.Warn().Err(err).Msg(httpServer)
				} else {
					l.Warn().Msg(httpServer)
				}
			default:
				l.Debug().Msg(httpServer)
			}

			return err
		}
	}
}
