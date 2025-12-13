package proxy

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/kamogelosekhukhune777/api-gateway/internal/observability/logger"
)

type Config struct {
	Log         *logger.Logger
	DialTimeout time.Duration // 5 * time.Second

	//The total time to wait for the response headers from the upstream.
	ResponseHeaderTimeout time.Duration //10 * time.Second

	//Keep-Alive period for idle connections in the pool
	KeepAlive time.Duration //30 * time.Second

	//maximum idle connections per host
	MaxIdleConnsPerHost int //100
}

// newHopHeaders defines headers that should not be proxied to the upstream server.
// These are often "hop-by-hop" headers relevant only to the connection between client and proxy.
var newHopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // Trailing encoding
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

// NewImprovedSingleHostReverseProxy creates an improved reverse proxy
func NewSingleHostReverseProxy(cfg Config, target string) http.Handler {
	ctx := context.Background()

	u, err := url.Parse(target)
	if err != nil {
		cfg.Log.Error(ctx, "invalid target URL format", "target", target, "err", err)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "bad gateway: invalid configuration", http.StatusBadGateway)
		})
	}

	rp := httputil.NewSingleHostReverseProxy(u)

	rp.Transport = &http.Transport{
		ResponseHeaderTimeout: cfg.ResponseHeaderTimeout,
		DialContext: (&net.Dialer{
			Timeout:   cfg.DialTimeout,
			KeepAlive: cfg.KeepAlive,
		}).DialContext,
		ForceAttemptHTTP2:   false,
		MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
	}

	originalDirector := rp.Director
	rp.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = u.Host

		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Gateway", "resilient-gateway")

		for _, h := range newHopHeaders {
			req.Header.Del(h)
		}
	}

	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			cfg.Log.Error(r.Context(), "proxy timeout error", "target", u.Host, "method", r.Method, "path", r.URL.Path, "err", err)
		} else {
			// Catch other connection/IO errors
			cfg.Log.Error(r.Context(), "proxy connection error", "target", u.Host, "method", r.Method, "path", r.URL.Path, "err", err)
		}

		http.Error(w, "upstream error", http.StatusBadGateway)
	}

	return rp
}
