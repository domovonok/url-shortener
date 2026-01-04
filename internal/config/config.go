package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type PoolConfig struct {
	MaxConnLifetime       time.Duration
	MaxConnLifetimeJitter time.Duration
	MaxConnIdleTime       time.Duration
	MaxConns              int32
	MinConns              int32
	MinIdleConns          int32
	HealthCheckPeriod     time.Duration
	PingMaxRetries        int
	PingRetryDelay        time.Duration
	PingTimeout           time.Duration
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Pool     PoolConfig
}

type ServerConfig struct {
	Port                    string
	GracefulShutdownTimeout time.Duration
	MetricsPeriod           time.Duration
}

type CacheConfig struct {
	Host           string
	Port           string
	Username       string
	Password       string
	DB             int
	PingTimeout    time.Duration
	PingMaxRetries int
	PingRetryDelay time.Duration
	Ttl            time.Duration
}

type Config struct {
	Server ServerConfig
	DB     DBConfig
	Cache  CacheConfig
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		Server: ServerConfig{
			Port:                    getEnvAsString("PORT", "8080"),
			GracefulShutdownTimeout: getEnvAsDuration("GRACEFUL_SHUTDOWN_TIMEOUT", 5*time.Second),
			MetricsPeriod:           getEnvAsDuration("METRICS_PERIOD", 5*time.Second),
		},
		DB: DBConfig{
			Host:     getEnvAsString("POSTGRES_HOST", "localhost"),
			Port:     getEnvAsString("POSTGRES_PORT", "5432"),
			Name:     getEnvAsString("POSTGRES_DB", "testdb"),
			User:     getEnvAsString("POSTGRES_USER", "myuser"),
			Password: getEnvAsString("POSTGRES_PASSWORD", "mypassword"),
			Pool: PoolConfig{
				MaxConnLifetime:       getEnvAsDuration("POSTGRES_MAX_CONN_LIFETIME", time.Hour),
				MaxConnLifetimeJitter: getEnvAsDuration("POSTGRES_MAX_CONN_LIFETIME_JITTER", 5*time.Minute),
				MaxConnIdleTime:       getEnvAsDuration("POSTGRES_MAX_CONN_IDLE_TIME", 30*time.Minute),
				MaxConns:              getEnvAsInt32("POSTGRES_MAX_CONNS", 20),
				MinConns:              getEnvAsInt32("POSTGRES_MIN_CONNS", 5),
				MinIdleConns:          getEnvAsInt32("POSTGRES_MIN_IDLE_CONNS", 2),
				HealthCheckPeriod:     getEnvAsDuration("POSTGRES_HEALTH_CHECK_PERIOD", time.Minute),
				PingMaxRetries:        getEnvAsInt("POSTGRES_MAX_RETRIES", 5),
				PingRetryDelay:        getEnvAsDuration("POSTGRES_RETRY_DELAY", time.Second),
				PingTimeout:           getEnvAsDuration("POSTGRES_AWAIT_TIME", 10*time.Second),
			},
		},
		Cache: CacheConfig{
			Host:           getEnvAsString("REDIS_HOST", "localhost"),
			Port:           getEnvAsString("REDIS_PORT", "6379"),
			Username:       getEnvAsString("REDIS_USERNAME", ""),
			Password:       getEnvAsString("REDIS_PASSWORD", ""),
			DB:             getEnvAsInt("REDIS_DB", 0),
			PingTimeout:    getEnvAsDuration("REDIS_PING_TIMEOUT", 5*time.Second),
			PingMaxRetries: getEnvAsInt("REDIS_PING_MAX_RETRIES", 3),
			PingRetryDelay: getEnvAsDuration("REDIS_PING_RETRY_DELAY", time.Second),
			Ttl:            getEnvAsDuration("CACHE_TTL", 10*time.Minute),
		},
	}
}

func getEnvAs[T any](key string, defaultVal T, parse func(string) (T, error)) T {
	if val := os.Getenv(key); val != "" {
		if v, err := parse(val); err == nil {
			return v
		}
	}
	return defaultVal
}

func getEnvAsString(key string, defaultVal string) string {
	return getEnvAs[string](key, defaultVal, func(s string) (string, error) {
		return s, nil
	})
}
func getEnvAsInt(key string, defaultVal int) int {
	return getEnvAs[int](key, defaultVal, strconv.Atoi)
}

func getEnvAsInt32(key string, defaultVal int32) int32 {
	return getEnvAs[int32](key, defaultVal, func(s string) (int32, error) {
		v, err := strconv.ParseInt(s, 10, 32)
		return int32(v), err
	})
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	return getEnvAs(key, defaultVal, time.ParseDuration)
}
