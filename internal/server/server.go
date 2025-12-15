package server

import (
	"log"
	"net/http"

	"github.com/kamogelosekhukhune777/api-gateway/internal/middleware"
	"github.com/kamogelosekhukhune777/api-gateway/internal/router"
)

type ServerConfig struct {
	RouterConfig router.Config
}

// NewHTTPServer configures the complete HTTP server handler,
// explicitly setting the global middleware chain around the router.
func NewServer(log *log.Logger, cfg *ServerConfig) http.Handler {
	h := router.NewRouter(cfg.RouterConfig)

	h = middleware.Logger(cfg.RouterConfig.Log, h)

	return h
}
