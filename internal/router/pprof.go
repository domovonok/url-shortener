package router

import (
	"net/http"

	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
)

func NewPprofRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/", http.DefaultServeMux)
	return r
}
