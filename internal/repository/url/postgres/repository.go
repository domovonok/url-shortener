package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	xerrors "github.com/domovonok/url-shortener/internal/errors"
	urlsqlc "github.com/domovonok/url-shortener/internal/repository/url/postgres/sqlc"
)

type repo struct {
	q *urlsqlc.Queries
}

func NewRepository(db urlsqlc.DBTX) *repo {
	return &repo{q: urlsqlc.New(db)}
}

func (r *repo) Create(ctx context.Context, url string) (int64, error) {
	id, err := r.q.Create(ctx, url)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, xerrors.ErrUrlExists
	}

	return id, err
}

func (r *repo) Get(ctx context.Context, id int64) (string, error) {
	url, err := r.q.Get(ctx, id)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", xerrors.ErrUrlNotFound
	}

	return url, err
}

func (r *repo) GetIDByURL(ctx context.Context, url string) (int64, error) {
	id, err := r.q.GetIDByURL(ctx, url)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, xerrors.ErrUrlNotFound
	}

	return id, err
}
