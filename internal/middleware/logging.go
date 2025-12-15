package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kamogelosekhukhune777/api-gateway/internal/observability/logger"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(b)
}

func Logger(log *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()

		path := r.URL.Path
		if r.URL.RawQuery != "" {
			path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
		}

		log.Info(r.Context(), "request started", "method", r.Method, "path", path, "remoteaddr", r.RemoteAddr)

		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		log.Info(r.Context(), "request completed", "method", r.Method, "path", path, "remoteaddr", r.RemoteAddr,
			"statuscode", recorder.status, "since", time.Since(now).String())

	})
}
