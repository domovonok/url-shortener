package router

import (
	"github.com/domovonok/url-shortener/internal/transport/http/common"
	"github.com/go-chi/chi/v5"
)

func New(linkHandler LinkHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Head("/healthcheck", common.Healthcheck)
	r.Post("/", linkHandler.Create)
	r.Get("/{code}", linkHandler.Get)

	return r
}
