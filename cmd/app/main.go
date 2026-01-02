package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/domovonok/url-shortener/internal/config"
	"github.com/domovonok/url-shortener/internal/database"
	"github.com/domovonok/url-shortener/internal/logger"
	linkRepo "github.com/domovonok/url-shortener/internal/repo/link"
	"github.com/domovonok/url-shortener/internal/router"
	linkHandler "github.com/domovonok/url-shortener/internal/transport/http/link"
	linkCreateUsecase "github.com/domovonok/url-shortener/internal/usecase/link/create"
	linkGetUsecase "github.com/domovonok/url-shortener/internal/usecase/link/get"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	log := logger.NewZapLogger(zapLogger)
	defer log.Sync()

	dbPool := database.MustInit(cfg.DB, log)
	repo := linkRepo.New(dbPool)

	startServer(
		dbPool,
		linkHandler.New(
			linkCreateUsecase.New(repo),
			linkGetUsecase.New(repo),
			log),
		cfg.Server,
		log,
	)
}

func startServer(
	dbPool *pgxpool.Pool,
	linkHandler *linkHandler.Controller,
	cfg config.ServerConfig,
	log logger.Logger,
) {
	srv := &http.Server{
		Addr:    net.JoinHostPort("", cfg.Port),
		Handler: router.New(linkHandler),
	}

	serverErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()
	log.Info("Server listening on", logger.Any("addr", srv.Addr))

	waitGracefulShutdown(srv, dbPool, serverErr, cfg.GracefulShutdownTimeout, log)

	log.Info("Service stopped successfully")
}

func waitGracefulShutdown(srv *http.Server, dbPool *pgxpool.Pool, serverErr <-chan error, timeout time.Duration, log logger.Logger) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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

	log.Info("Closing DB pool...")
	dbPool.Close()
	log.Info("DB pool closed")
}
