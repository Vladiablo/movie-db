package main

import (
	"context"
	"log/slog"
	"movie-db/internal"
	"movie-db/internal/http"
	"movie-db/internal/http/handlers/movies"
	"movie-db/internal/repository"
	"movie-db/internal/service"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type args struct {
	DbUrl                   string        `arg:"--db-url,env:DATABASE_URL" help:"database URL"`
	DbTimeout               time.Duration `arg:"--db-timeout,env:DATABASE_TIMEOUT" default:"10s" validate:"min=100ms" help:"database requests timeout"`
	GracefulShutdownTimeout time.Duration `arg:"--graceful-shutdown-timeout" default:"10s" validate:"min=0s" help:"graceful shutdown timeout"`
	HttpTimeout             time.Duration `arg:"--http-timeout" default:"30s" validate:"min=100ms" help:"http requests timeout"`
}

//	@title			Movie DB API
//	@version		1.0
//	@description	API to get movies and series information

//	@host		localhost:8080
//	@BasePath	/v1

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

func run() int {
	_ = godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	validator := internal.NewValidator()

	var args args
	if err := arg.Parse(&args); err != nil {
		logger.Error("Failed to parse program arguments. Exiting...", slog.Any("error", err))
		return 1
	}

	if err := validator.Validate(&args); err != nil {
		logger.Error("Failed to validate program arguments. Exiting...", slog.Any("error", err))
		return 1
	}

	cfg, err := pgxpool.ParseConfig(args.DbUrl)
	if err != nil {
		logger.Error("Failed to parse database URL. Exiting...", slog.Any("error", err))
		return 1
	}

	var pool *pgxpool.Pool
	{
		ctx, cancel := context.WithTimeout(context.Background(), args.GracefulShutdownTimeout)
		defer cancel()

		pool, err = pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			logger.Error("Failed to create database connection pool. Exiting...", slog.Any("error", err))
			return 1
		}
		defer pool.Close()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := repository.NewPgRepository(pool, args.DbTimeout)

	svc := service.New(repo,
		logger.With(slog.String("component", "service")),
	)

	srv := http.NewServer(&http.ServerConfig{
		Timeout:                 args.HttpTimeout,
		GracefulShutdownTimeout: args.GracefulShutdownTimeout,
		Validator:               validator,
		Logger:                  logger.With(slog.String("component", "http-server")),
	})

	moviesHandler := movies.New(svc)
	http.RegisterRoutes(srv, moviesHandler)

	var wg sync.WaitGroup
	wg.Add(1)

	var firstErr atomic.Pointer[error]

	go func() {
		defer wg.Done()
		defer cancel()

		if err := srv.Start(ctx); err != nil {
			firstErr.CompareAndSwap(nil, &err)

			logger.Error("Failed to serve HTTP", slog.Any("error", err))
		}
	}()

	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)

		wg.Wait()
	}()

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	var timeoutCh <-chan time.Time

	ctxDoneCh := ctx.Done()

	for {
		select {
		case sig := <-sigCh:
			logger.Info("Received signal. Shutting down...", slog.String("signal", sig.String()))

			cancel()

		case <-ctxDoneCh:
			if args.GracefulShutdownTimeout == 0 {
				return 0
			}

			ctxDoneCh = nil
			timeoutCh = time.After(args.GracefulShutdownTimeout)

		case <-timeoutCh:
			logger.Error("Failed to graceful shutdown", slog.Duration("duration", args.GracefulShutdownTimeout))

			return 1

		case <-doneCh:
			if firstErr.Load() != nil {
				return 1
			}

			logger.Info("Graceful shutdown complete")
			return 0
		}
	}
}

func main() {
	os.Exit(run())
}
