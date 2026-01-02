package get

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

func (s *Usecase) Get(ctx context.Context, code string) (model.Link, error) {
	return s.link.Get(ctx, code)
}
