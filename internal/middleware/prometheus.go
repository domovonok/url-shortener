package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/domovonok/url-shortener/internal/metrics"
)

func Prometheus(m *metrics.PrometheusMetrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r)

			duration := time.Since(start).Seconds()

			routePattern := chi.RouteContext(r.Context()).RoutePattern()
			if routePattern == "" {
				routePattern = r.URL.Path
			}

			m.HTTPRequestsTotal.WithLabelValues(r.Method, routePattern, strconv.Itoa(rw.statusCode)).Inc()
			m.HTTPRequestDuration.WithLabelValues(routePattern).Observe(duration)
		})
	}
}
