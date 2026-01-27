package testdb

import (
	"context"
	"database/sql"
	"os"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	syncOnce  sync.Once
	connStr   string
	container *postgres.PostgresContainer
)

func TestWithMigrations() (*postgres.PostgresContainer, string, error) {
	ctx := context.Background()
	var creationError error

	syncOnce.Do(func() {
		container, creationError = postgres.Run(ctx,
			"postgres:15-alpine",
			postgres.WithDatabase(os.Getenv("POSTGRES_DB")),
			postgres.WithUsername(os.Getenv("POSTGRES_USER")),
			postgres.WithPassword(os.Getenv("POSTGRES_PASSWORD")),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(5*time.Minute)),
		)
		if creationError != nil {
			return
		}

		connStr, creationError = container.ConnectionString(ctx, "sslmode=disable")
		if creationError != nil {
			return
		}

		db, err := sql.Open("pgx", connStr)
		if err != nil {
			creationError = err
			return
		}
		defer func() {
			_ = db.Close()
		}()

		if err := db.PingContext(ctx); err != nil {
			creationError = err
			return
		}

		if err := goose.Up(db, "../../migrations"); err != nil {
			creationError = err
			return
		}
	})

	if creationError != nil {
		return nil, "", creationError
	}

	return container, connStr, nil
}
