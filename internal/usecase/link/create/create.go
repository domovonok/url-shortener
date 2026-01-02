package create

import (
	"context"

	"github.com/domovonok/url-shortener/internal/model"
)

type Usecase struct {
	link linkRepo
}

func New(l linkRepo) *Usecase {
	return &Usecase{link: l}
}

func (s *Usecase) Create(ctx context.Context, url string) (model.Link, error) {
	return s.link.Create(ctx, url)
}
