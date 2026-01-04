package database

import (
	"context"
	"fmt"
	"time"

	"github.com/domovonok/url-shortener/internal/config"
	"github.com/domovonok/url-shortener/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MustInit(cfg config.DBConfig, log logger.Logger) *pgxpool.Pool {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatal("invalid connection string: %w", logger.Error(err))
	}

	poolConfig.MaxConnLifetime = cfg.Pool.MaxConnLifetime
	poolConfig.MaxConnLifetimeJitter = cfg.Pool.MaxConnLifetimeJitter
	poolConfig.MaxConnIdleTime = cfg.Pool.MaxConnIdleTime
	poolConfig.MaxConns = cfg.Pool.MaxConns
	poolConfig.MinConns = cfg.Pool.MinConns
	poolConfig.MinIdleConns = cfg.Pool.MinIdleConns
	poolConfig.HealthCheckPeriod = cfg.Pool.HealthCheckPeriod

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Pool.PingTimeout)
	defer cancel()

	dbPool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatal("Failed to initialize database:", logger.Error(err))
	}

	var pingErr error
	for i := 1; i <= cfg.Pool.PingMaxRetries; i++ {
		pingCtx, pingCancel := context.WithTimeout(context.Background(), cfg.Pool.PingTimeout)
		pingErr = dbPool.Ping(pingCtx)
		pingCancel()
		if pingErr == nil {
			break
		}
		log.Warn(
			"Database ping attempt failed",
			logger.Any("attempt", i),
			logger.Any("max_retries", cfg.Pool.PingMaxRetries),
			logger.Error(pingErr),
		)
		if i < cfg.Pool.PingMaxRetries {
			log.Warn(
				"Retrying database ping",
				logger.Any("delay", cfg.Pool.PingRetryDelay),
			)
			time.Sleep(cfg.Pool.PingRetryDelay)
		}
	}

	if pingErr != nil {
		dbPool.Close()
		log.Fatal("Unable to ping database")
	}

	log.Info("Database connection established")
	return dbPool
}
