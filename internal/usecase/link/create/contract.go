//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_test -destination mocks_test.go
package create

import (
	"context"

	"github.com/domovonok/url-shortener/internal/model"
)

type linkRepo interface {
	Create(ctx context.Context, url string) (model.Link, error)
}
