package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/domovonok/url-shortener/internal/config"
	urlhandler "github.com/domovonok/url-shortener/internal/handler/url"
	postgresinfra "github.com/domovonok/url-shortener/internal/infrastructure/postgres"
	urlinmemoryrepo "github.com/domovonok/url-shortener/internal/repository/url/inmemory"
	urlpostgresrepo "github.com/domovonok/url-shortener/internal/repository/url/postgres"
	"github.com/domovonok/url-shortener/internal/router"
	urlservice "github.com/domovonok/url-shortener/internal/service/url"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalln("can not initialize config:", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("can not initialize logger:", err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Println("failed to sync logger:", err)
		}
	}(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var repo urlservice.Repo

	switch cfg.DB.Storage {
	case config.StorageInMemory:
		repo = urlinmemoryrepo.NewRepository()
	case config.StoragePostgres:
		pool, err := pgxpool.New(ctx, cfg.DB.Postgres.DSN())

		if err != nil {
			logger.Fatal("failed to connect to postgres", zap.Error(err))
		}
		defer pool.Close()

		postgresinfra.SetupPostgres(pool, logger)

		repo = urlpostgresrepo.NewRepository(pool)
	default:
		logger.Fatal("unknown storage type", zap.String("storage", string(cfg.DB.Storage)))
	}

	service := urlservice.NewService(repo)

	handler := urlhandler.NewHandler(logger, service)

	runServer(ctx, &cfg.Server, logger, handler)
}

func runServer(ctx context.Context, cfg *config.ServerConfig, logger *zap.Logger, handler router.Handler) {
	srv := &http.Server{
		Addr:    cfg.Addr(),
		Handler: router.New(handler),
	}

	serverErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()
	logger.Info("Server listening on", zap.String("addr", srv.Addr))

	waitGracefulShutdown(ctx, logger, srv, serverErr, cfg.GracefulShutdownTimeout)

	logger.Info("Service stopped successfully")
}

func waitGracefulShutdown(
	ctx context.Context,
	logger *zap.Logger,
	srv *http.Server,
	serverErr <-chan error,
	timeout time.Duration,
) {
	var reason string
	select {
	case <-ctx.Done():
		reason = "signal"
	case err := <-serverErr:
		reason = "server error: " + err.Error()
	}

	logger.Info("Shutting down...", zap.String("reason", reason))

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error", zap.Error(err))
	} else {
		logger.Info("HTTP server stopped")
	}
}
