package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kamogelosekhukhune777/api-gateway/internal/trace"
)

const TraceHeader = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := uuid.New()

		traceID := r.Header.Get(TraceHeader)
		if traceID == "" {
			traceID = tracer.String()
		}

		ctx := trace.SetTraceID(r.Context(), tracer)

		w.Header().Set(TraceHeader, traceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
