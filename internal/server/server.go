package server

import (
	"net/http"

	"github.com/kamogelosekhukhune777/api-gateway/internal/router"
)

type ServerConfig struct {
	RouterConfig router.Config
}

// NewHTTPServer configures the complete HTTP server handler,
// explicitly setting the global middleware chain around the router.
func NewHTTPServer(cfg *ServerConfig) http.Handler {
	r := router.NewRouter(cfg.RouterConfig)

	h := r

	return h
}
