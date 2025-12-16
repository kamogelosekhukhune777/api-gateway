package middleware

import (
	"net/http"

	"github.com/kamogelosekhukhune777/api-gateway/internal/observability/metrics"
)

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Inject metrics into context (singleton pointer)
		ctx := metrics.Set(r.Context())

		// Count total requests
		metrics.AddRequests(ctx)

		// Track goroutine count opportunistically
		metrics.AddGoroutines(ctx)

		// Execute downstream handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
