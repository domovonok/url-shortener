package middleware

import (
	"net/http"
	"time"

	"github.com/domovonok/url-shortener/internal/logger"
	"github.com/go-chi/chi/v5"
)

func Logger(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			routePattern := chi.RouteContext(r.Context()).RoutePattern()
			if routePattern == "" {
				routePattern = r.URL.Path
			}

			log.Debug("HTTP request",
				logger.Any("method", r.Method),
				logger.Any("path", routePattern),
				logger.Any("status", rw.statusCode),
				logger.Any("duration", duration),
			)
		})
	}
}
