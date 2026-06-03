package config

import (
	"net"
	"net/url"
	"time"

	"github.com/caarlos0/env/v11"
)

type StorageType string

const (
	StorageInMemory StorageType = "inmemory"
	StoragePostgres StorageType = "postgres"
)

type (
	Config struct {
		Debug  bool `env:"DEBUG" envDefault:"false"`
		Server ServerConfig
		DB     DBConfig
	}

	ServerConfig struct {
		Addr                    string        `env:"HTTP_ADDR" envDefault:":8080"`
		GracefulShutdownTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT" envDefault:"5s"`
	}

	DBConfig struct {
		Storage  StorageType `env:"STORAGE_TYPE" envDefault:"inmemory"`
		Postgres PostgresConfig
	}

	PostgresConfig struct {
		Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
		Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
		DB       string `env:"POSTGRES_DB" envDefault:"testdb"`
		User     string `env:"POSTGRES_USER" envDefault:"user"`
		Password string `env:"POSTGRES_PASSWORD" envDefault:"12345"`
		SSLMode  string `env:"POSTGRES_SSL_MODE" envDefault:"disable"`
	}
)

func (c *PostgresConfig) DSN() string {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, c.Port),
		Path:   c.DB,
	}

	q := dsn.Query()
	q.Set("sslmode", c.SSLMode)
	dsn.RawQuery = q.Encode()

	return dsn.String()
}

func New() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return &cfg, err
}
