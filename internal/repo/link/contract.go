package link

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type dbPool interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
