package get

import (
	"context"

	"github.com/domovonok/url-shortener/internal/model"
)

type linkRepo interface {
	Get(ctx context.Context, code string) (model.Link, error)
}
