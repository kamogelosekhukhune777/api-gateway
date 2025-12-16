package server

import (
	"net/http"

	"github.com/kamogelosekhukhune777/api-gateway/internal/middleware"
	"github.com/kamogelosekhukhune777/api-gateway/internal/router"
)

type Config struct {
	RouterConfig router.Config
}

// NewHTTPServer configures the complete HTTP server handler,
// explicitly setting the global middleware chain around the router.
func NewServer(cfg *Config) http.Handler {
	h := router.NewRouter(cfg.RouterConfig)

	h = middleware.Panics(cfg.RouterConfig.Log, h)
	h = middleware.Metrics(h)
	h = middleware.Logger(cfg.RouterConfig.Log, h)
	h = middleware.RequestID(h)

	return h
}
