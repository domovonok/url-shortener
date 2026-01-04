package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/domovonok/url-shortener/internal/cache"
	"github.com/domovonok/url-shortener/internal/config"
	"github.com/domovonok/url-shortener/internal/database"
	"github.com/domovonok/url-shortener/internal/logger"
	"github.com/domovonok/url-shortener/internal/metrics"
	linkRepo "github.com/domovonok/url-shortener/internal/repo/link"
	"github.com/domovonok/url-shortener/internal/router"
	linkHandler "github.com/domovonok/url-shortener/internal/transport/http/link"
	linkCreateUsecase "github.com/domovonok/url-shortener/internal/usecase/link/create"
	linkGetUsecase "github.com/domovonok/url-shortener/internal/usecase/link/get"
)

func main() {
	cfg := config.Load()

	log := logger.MustInit(cfg.Debug)
	defer func() {
		if err := log.Sync(); err != nil {
			log.Error("Unable to sync logger", logger.Error(err))
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbPool := database.MustInit(cfg.DB, log)
	defer dbPool.Close()
	repo := linkRepo.New(dbPool)

	dbCache := cache.MustInit(cfg.Cache, log)
	defer func() {
		if err := dbCache.Close(); err != nil {
			log.Error("Unable to close cache", logger.Error(err))
		}
	}()

	cacheRepo := linkRepo.NewCached(repo, dbCache, log)

	startServer(
		ctx,
		linkHandler.New(
			linkCreateUsecase.New(cacheRepo),
			linkGetUsecase.New(cacheRepo),
			log),
		cfg.Server,
		log,
	)
}

func startServer(
	ctx context.Context,
	linkHandler *linkHandler.Controller,
	cfg config.ServerConfig,
	log logger.Logger,
) {
	prom := metrics.NewPrometheusMetrics()
	metrics.StartSystemMetricsCollector(ctx, prom, cfg.MetricsPeriod)

	srv := &http.Server{
		Addr:    net.JoinHostPort("", cfg.Port),
		Handler: router.New(linkHandler, log, prom),
	}

	serverErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()
	log.Info("Server listening on", logger.Any("addr", srv.Addr))

	waitGracefulShutdown(ctx, srv, serverErr, cfg.GracefulShutdownTimeout, log)

	log.Info("Service stopped successfully")
}

func waitGracefulShutdown(ctx context.Context, srv *http.Server, serverErr <-chan error, timeout time.Duration, log logger.Logger) {
	var reason string
	select {
	case <-ctx.Done():
		reason = "signal"
	case err := <-serverErr:
		reason = "server error: " + err.Error()
	}

	log.Info("Shutting down...", logger.Any("reason", reason))
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("HTTP server shutdown error", logger.Error(err))
	} else {
		log.Info("HTTP server stopped")
	}
}
