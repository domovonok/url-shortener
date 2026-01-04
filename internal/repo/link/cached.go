package link

import (
	"context"
	"encoding/json"

	"github.com/domovonok/url-shortener/internal/logger"
	"github.com/domovonok/url-shortener/internal/model"
)

type CachedRepo struct {
	r   BaseRepo
	c   Cache
	log logger.Logger
}

func NewCached(r BaseRepo, c Cache, l logger.Logger) *CachedRepo {
	return &CachedRepo{r: r, c: c, log: l}
}

func (cr *CachedRepo) Create(ctx context.Context, url string) (model.Link, error) {
	res, err := cr.r.Create(ctx, url)
	if err == nil {
		if data, err := json.Marshal(res); err == nil {
			_ = cr.c.Set(ctx, key(res.Code), data)
		}
	}
	return res, err
}

func (cr *CachedRepo) Get(ctx context.Context, code string) (model.Link, error) {
	if data, err := cr.c.Get(ctx, key(code)); err == nil {
		var l model.Link
		if json.Unmarshal(data, &l) == nil {
			cr.log.Debug("Cache hit", logger.Any("code", code))
			return l, nil
		}
	} else {
		cr.log.Debug("Cache miss", logger.Any("code", code), logger.Error(err))
	}

	res, err := cr.r.Get(ctx, code)
	if err != nil {
		return model.Link{}, err
	}

	if data, err := json.Marshal(res); err == nil {
		_ = cr.c.Set(ctx, key(code), data)
	}

	return cr.r.Get(ctx, code)
}

func key(code string) string {
	return "link:" + code
}
