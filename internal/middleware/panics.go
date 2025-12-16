package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/kamogelosekhukhune777/api-gateway/internal/observability/logger"
	"github.com/kamogelosekhukhune777/api-gateway/internal/observability/metrics"
)

func Panics(log *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if rec := recover(); rec != nil {
				trace := debug.Stack()

				// Log panic with trace + trace_id already in context
				log.Error(r.Context(), "panic recovered", "panic", rec, "stack", string(trace))

				// Increment panic metric
				metrics.AddPanics(r.Context())

				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
