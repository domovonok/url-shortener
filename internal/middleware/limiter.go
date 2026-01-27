package middleware

import (
	"net/http"
	"strconv"

	"github.com/domovonok/url-shortener/internal/logger"
	"github.com/domovonok/url-shortener/internal/metrics"
)

type tokenBucket interface {
	Allow() bool
	Capacity() int
	Remaining() int
}

func RateLimitMiddleware(limiter tokenBucket, log logger.Logger, m *metrics.PrometheusMetrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				m.RateLimitExceededTotal.Inc()
				log.Warn("Rate limit exceeded",
					logger.Any("method", r.Method),
					logger.Any("path", r.URL.Path),
					logger.Any("remote_addr", r.RemoteAddr),
				)

				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limiter.Capacity()))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte("Rate limit exceeded"))
				return
			}

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limiter.Capacity()))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(limiter.Remaining()))

			next.ServeHTTP(w, r)
		})
	}
}
