package link

import (
	"context"

	"github.com/domovonok/url-shortener/internal/model"
)

type createUsecase interface {
	Create(ctx context.Context, url string) (model.Link, error)
}

type getUsecase interface {
	Get(ctx context.Context, code string) (model.Link, error)
}
