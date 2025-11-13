package startup

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	emiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"

	"github.com/kxddry/avito-backend-internship-2025/internal/api"
	"github.com/kxddry/avito-backend-internship-2025/internal/api/generated"
	"github.com/kxddry/avito-backend-internship-2025/internal/domain"
	httpmiddleware "github.com/kxddry/avito-backend-internship-2025/pkg/middleware"
)

// Application is the application.
type Application struct {
	cfg    *Config
	echo   *echo.Echo
	server *http.Server
}

// NewApplication creates a new application.
func NewApplication(cfg *Config, service domain.AssignmentService) (*Application, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(emiddleware.RequestID())
	e.Use(emiddleware.Recover())
	e.Use(httpmiddleware.Zerolog())

	handler := api.NewServer(service)
	strict := generated.NewStrictHandler(handler, nil)
	generated.RegisterHandlers(e, strict)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.ServerConfig.Port),
		Handler:           e,
		ReadHeaderTimeout: cfg.ServerConfig.Timeout,
		ReadTimeout:       cfg.ServerConfig.Timeout,
		WriteTimeout:      cfg.ServerConfig.Timeout,
		IdleTimeout:       cfg.ServerConfig.IdleTimeout,
	}

	return &Application{
		cfg:    cfg,
		echo:   e,
		server: httpServer,
	}, nil
}

// Run runs the application.
func (a *Application) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		log.Info().Msgf("HTTP server is listening on %s", a.server.Addr)
		if err := a.echo.StartServer(a.server); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout(a.cfg))
		defer cancel()
		return a.echo.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

// defTimeout is the default timeout.
const defTimeout = 5 * time.Second

// shutdownTimeout returns the shutdown timeout.
func shutdownTimeout(cfg *Config) time.Duration {
	timeout := cfg.ServerConfig.Timeout
	if timeout <= 0 {
		return defTimeout
	}
	return timeout
}
