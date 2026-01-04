package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/domovonok/url-shortener/internal/config"
	"github.com/domovonok/url-shortener/internal/logger"
)

type RedisCache struct {
	c   *redis.Client
	ttl time.Duration
}

func MustInit(cfg config.CacheConfig, log logger.Logger) *RedisCache {
	opts := &redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	}
	c := redis.NewClient(opts)

	redis.SetLogger(logger.NewRedisLogger(log))

	var pingErr error
	for i := 1; i <= cfg.PingMaxRetries; i++ {
		pingCtx, pingCancel := context.WithTimeout(context.Background(), cfg.PingTimeout)
		pingErr = c.Ping(pingCtx).Err()
		pingCancel()
		if pingErr == nil {
			break
		}
		log.Warn(
			"Redis ping attempt failed",
			logger.Any("attempt", i),
			logger.Any("max_retries", cfg.PingMaxRetries),
			logger.Error(pingErr),
		)
		if i < cfg.PingMaxRetries {
			log.Warn(
				"Retrying redis ping",
				logger.Any("delay", cfg.PingRetryDelay),
			)
			time.Sleep(cfg.PingRetryDelay)
		}
	}

	if pingErr != nil {
		_ = c.Close()
		log.Fatal("Unable to ping redis", logger.Error(pingErr))
	}

	log.Info("Redis connection established", logger.Any("addr", opts.Addr))
	return &RedisCache{c: c, ttl: cfg.Ttl}
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	return r.c.Get(ctx, key).Bytes()
}

func (r *RedisCache) Set(ctx context.Context, key string, value []byte) error {
	return r.c.Set(ctx, key, value, r.ttl).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.c.Del(ctx, key).Err()
}

func (r *RedisCache) Close() error {
	return r.c.Close()
}
