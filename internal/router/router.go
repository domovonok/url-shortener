package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}

func New(handler Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", handler.Create)
	r.Get("/{code}", handler.Get)

	return r
}
