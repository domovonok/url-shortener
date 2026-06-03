package inmemory

import (
	"context"
	"sync"

	"github.com/domovonok/url-shortener/internal/errors"
)

type repo struct {
	mu      sync.Mutex
	idToUrl map[int64]string
	urlToId map[string]int64
}

func NewRepository() *repo {
	return &repo{
		idToUrl: make(map[int64]string),
		urlToId: make(map[string]int64),
	}
}

func (r *repo) Create(_ context.Context, url string) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if id, ok := r.urlToId[url]; ok {
		return id, nil
	}

	id := int64(len(r.idToUrl) + 1)
	r.idToUrl[id] = url
	r.urlToId[url] = id

	return id, nil
}

func (r *repo) Get(_ context.Context, id int64) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	url, ok := r.idToUrl[id]
	if !ok {
		return "", errors.ErrUrlNotFound
	}

	return url, nil
}
