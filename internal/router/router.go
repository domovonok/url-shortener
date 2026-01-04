package router

import (
	"github.com/domovonok/url-shortener/internal/logger"
	"github.com/domovonok/url-shortener/internal/metrics"
	"github.com/domovonok/url-shortener/internal/middleware"
	"github.com/domovonok/url-shortener/internal/transport/http/common"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(linkHandler LinkHandler, log logger.Logger, prom *metrics.PrometheusMetrics) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Prometheus(prom))
	r.Use(middleware.Logger(log))

	r.Handle("/metrics", promhttp.Handler())

	r.Head("/healthcheck", common.Healthcheck)
	r.Post("/", linkHandler.Create)
	r.Get("/{code}", linkHandler.Get)

	return r
}
