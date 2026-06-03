package postgres

import (
	"database/sql"
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func SetupPostgres(pool *pgxpool.Pool, logger *zap.Logger) {
	goose.SetBaseFS(embedMigrations)
	goose.SetTableName("url_goose_db_version")

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Fatal("can not set dialect in goose", zap.Error(err))
	}

	db := stdlib.OpenDBFromPool(pool)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error("can not close db", zap.Error(err))
		}
	}(db)

	if err := goose.Up(db, "migrations"); err != nil {
		logger.Fatal("can not setup migrations", zap.Error(err))
	}
}
