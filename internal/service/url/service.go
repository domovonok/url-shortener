package url

import (
	"context"
	"errors"

	xerrors "github.com/domovonok/url-shortener/internal/errors"
)

type Repo interface {
	Create(ctx context.Context, url string) (int64, error)
	GetIDByURL(ctx context.Context, url string) (int64, error)
	Get(ctx context.Context, id int64) (string, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) *service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, url string) (string, error) {
	if err := validateURL(url); err != nil {
		return "", err
	}

	id, err := s.repo.Create(ctx, url)
	if errors.Is(err, xerrors.ErrUrlExists) {
		id, err = s.repo.GetIDByURL(ctx, url)
	}
	if err != nil {
		return "", err
	}

	code, err := idToCode(id)
	if err != nil {
		return "", errors.New("unreachable, code should be valid")
	}

	return code, nil
}

func (s *service) Get(ctx context.Context, code string) (string, error) {
	id, err := codeToID(code)
	if err != nil {
		return "", err
	}

	return s.repo.Get(ctx, id)
}
