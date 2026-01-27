package link

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/domovonok/url-shortener/internal/model"
)

type dbPool interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type cache interface {
	Set(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
}

type baseRepo interface {
	Get(ctx context.Context, code string) (model.Link, error)
	Create(ctx context.Context, url string) (model.Link, error)
}
