package link

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/domovonok/url-shortener/internal/model"
	"github.com/domovonok/url-shortener/internal/repo/link/codec"
	"github.com/jackc/pgx/v5"
)

const tableLinks = "links"

type Repo struct {
	pool         dbPool
	queryBuilder sq.StatementBuilderType
}

func New(pool dbPool) *Repo {
	return &Repo{
		pool:         pool,
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Repo) Create(ctx context.Context, url string) (model.Link, error) {
	query, args, _ := r.queryBuilder.
		Insert(tableLinks).
		Columns("url").
		Values(url).
		Suffix("ON CONFLICT (url) DO UPDATE SET url = EXCLUDED.url RETURNING id, created_at").
		ToSql()

	var (
		id        int64
		createdAt time.Time
	)

	if err := r.pool.QueryRow(ctx, query, args...).Scan(&id, &createdAt); err != nil {
		return model.Link{}, handleDBError(err)
	}

	res := model.Link{
		Url:       url,
		Code:      codec.EncodeIDToCode(id),
		CreatedAt: createdAt,
	}

	return res, nil
}

func (r *Repo) Get(ctx context.Context, code string) (model.Link, error) {
	id, err := codec.DecodeCodeToID(code)
	if err != nil {
		return model.Link{}, model.ErrCodeNotFound
	}

	query, args, _ := r.queryBuilder.
		Select("url", "created_at").
		From(tableLinks).
		Where(sq.Eq{"id": id}).
		ToSql()

	var res model.Link
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&res.Url, &res.CreatedAt); err != nil {
		return model.Link{}, handleDBError(err)
	}

	return res, nil
}

func handleDBError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return model.ErrCodeNotFound
	}
	return err
}
