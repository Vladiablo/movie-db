package main

import (
	"context"
	"log/slog"
	"movie-db/internal/http"
	"movie-db/internal/http/handlers/movies"
	"movie-db/internal/repository"
	"movie-db/internal/service"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const (
	httpServerAddr = ":8080"
	dbConnTimeout  = 10 * time.Second
)

//	@title			Movie DB API
//	@version		1.0
//	@description	API to get movies and series information

//	@host		localhost:8080
//	@BasePath	/v1

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

func run() int {
	const shutdownTimeout = 10 * time.Second

	_ = godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	cfg, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error("Failed to parse database config. Exiting...", slog.Any("error", err))
		return 1
	}

	var pool *pgxpool.Pool
	{
		ctx, cancel := context.WithTimeout(context.Background(), dbConnTimeout)
		defer cancel()

		pool, err = pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			logger.Error("Failed to create database connection pool. Exiting...", slog.Any("error", err))
			return 2
		}
		defer pool.Close()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := repository.NewPgRepository(pool,
		logger.With(slog.String("component", "pg-repo")),
	)

	svc := service.New(repo)

	srv := http.NewServer(httpServerAddr,
		logger.With(slog.String("component", "http-server")),
	)

	handlersLogger := logger.With(slog.String("component", "http"))
	moviesHandler := movies.New(svc, handlersLogger)
	http.RegisterRoutes(srv, moviesHandler)

	var wg sync.WaitGroup
	wg.Add(1)

	var firstErr atomic.Pointer[error]

	go func() {
		defer wg.Done()
		defer cancel()

		if err := srv.Start(ctx, shutdownTimeout); err != nil {
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
			ctxDoneCh = nil
			timeoutCh = time.After(shutdownTimeout)

		case <-timeoutCh:
			logger.Error("Failed to graceful shutdown", slog.String("duration", shutdownTimeout.String()))

			return 1

		case <-doneCh:
			if firstErr.Load() != nil {
				return 3
			}

			logger.Info("Graceful shutdown complete")
			return 0
		}
	}
}

func main() {
	os.Exit(run())
}
