package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/domovonok/url-shortener/internal/cache"
	"github.com/domovonok/url-shortener/internal/config"
	"github.com/domovonok/url-shortener/internal/database"
	"github.com/domovonok/url-shortener/internal/limiter"
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

	rateLimiter := limiter.NewTokenBucket(cfg.RateLimit)

	prom := metrics.NewPrometheusMetrics()
	metrics.StartSystemMetricsCollector(ctx, prom, cfg.MetricsPeriod)

	startServer(
		ctx,
		linkHandler.New(
			linkCreateUsecase.New(cacheRepo),
			linkGetUsecase.New(cacheRepo),
			log),
		rateLimiter,
		prom,
		cfg.Server,
		log,
	)
}

func startServer(
	ctx context.Context,
	linkHandler router.LinkHandler,
	rateLimiter router.TokenBucket,
	prom *metrics.PrometheusMetrics,
	cfg config.ServerConfig,
	log logger.Logger,
) {
	mainSrv := &http.Server{
		Addr:    net.JoinHostPort("", cfg.Port),
		Handler: router.New(linkHandler, rateLimiter, log, prom),
	}

	serverErr := make(chan error, 1)
	go func() {
		if err := mainSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()
	log.Info("Server listening on", logger.Any("addr", mainSrv.Addr))

	pprofSrv := &http.Server{
		Addr:    net.JoinHostPort("", cfg.PprofPort),
		Handler: router.NewPprofRouter(),
	}

	go func() {
		if err := pprofSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Pprof server error", logger.Error(err))
		}
	}()
	log.Info("Pprof server listening on", logger.Any("addr", pprofSrv.Addr))

	waitGracefulShutdown(ctx, mainSrv, pprofSrv, serverErr, cfg.GracefulShutdownTimeout, log)

	log.Info("Service stopped successfully")
}

func waitGracefulShutdown(ctx context.Context, mainSrv, pprofSrv *http.Server, serverErr <-chan error, timeout time.Duration, log logger.Logger) {
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

	wg := new(sync.WaitGroup)

	wg.Go(func() {
		if err := pprofSrv.Shutdown(shutdownCtx); err != nil {
			log.Error("Pprof server graceful shutdown failed", logger.Error(err))
		} else {
			log.Info("Pprof server stopped")
		}
	})
	wg.Go(func() {
		if err := mainSrv.Shutdown(shutdownCtx); err != nil {
			log.Error("HTTP server shutdown error", logger.Error(err))
		} else {
			log.Info("HTTP server stopped")
		}
	})

	wg.Wait()
}
