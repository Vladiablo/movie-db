package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

const DefaultServerAddr = ":8080"

type Server struct {
	echo            *echo.Echo
	addr            string
	gracefulTimeout time.Duration
}

type ServerConfig struct {
	Addr                    string
	Timeout                 time.Duration
	GracefulShutdownTimeout time.Duration

	Validator echo.Validator
	Logger    *slog.Logger
}

func NewServer(cfg *ServerConfig) *Server {
	e := echo.NewWithConfig(echo.Config{
		Logger:           cfg.Logger,
		HTTPErrorHandler: httpErrorHandler,
	})

	addr := cfg.Addr
	if addr == "" {
		addr = DefaultServerAddr
	}

	e.Use(middleware.Recover())

	if cfg.Timeout > 0 {
		e.Use(middleware.ContextTimeout(cfg.Timeout))
	}

	e.Validator = cfg.Validator

	return &Server{
		echo:            e,
		addr:            addr,
		gracefulTimeout: max(0, cfg.GracefulShutdownTimeout),
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.echo.Logger.Info("Starting HTTP server...", slog.String("address", s.addr))

	// Note: lc.Listen doesn't respect cancelled context, so we do our own check
	if ctx.Err() != nil {
		return ctx.Err()
	}

	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, "tcp", s.addr)
	if err != nil {
		return fmt.Errorf("http server failed to listen: %w", err)
	}

	go func() {
		<-ctx.Done()

		s.echo.Logger.Info("Shutting down HTTP server...")
	}()

	sc := echo.StartConfig{
		Address:         s.addr,
		HideBanner:      true,
		Listener:        listener,
		GracefulTimeout: s.gracefulTimeout,
	}

	if err := sc.Start(ctx, s.echo); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			s.echo.Logger.Info("HTTP server stopped")

			return nil
		}

		s.echo.Logger.Error("Failed to start HTTP server", slog.Any("error", err))

		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}
