package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kamogelosekhukhune777/api-gateway/internal/observability/logger"
	"github.com/kamogelosekhukhune777/api-gateway/internal/proxy"
)

type Route struct {
	Prefix  string
	Service string
	Methods []string
}

type Config struct {
	Log         *logger.Logger
	Services    map[string]string
	Routes      []Route
	ProxyConfig proxy.Config
}

// NewRouter initializes the router, applies global middleware, and configures dynamic routes.
func NewRouter(cfg Config) http.Handler {
	m := mux.NewRouter()

	// Health Check Route
	m.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods("GET")
	cfg.Log.Info(context.Background(), "Registered static route", "path", "/health")

	//middleware chain here

	// DYNAMIC ROUTES FROM CONFIG
	for _, rt := range cfg.Routes {
		// Validation: Ensure we have a prefix
		if strings.TrimSpace(rt.Prefix) == "" {
			cfg.Log.Error(context.Background(), "Skipping route: Prefix is empty", "route_name", rt.Service)
			continue
		}

		svcAddr, ok := cfg.Services[rt.Service]
		if !ok {
			cfg.Log.Error(context.Background(), "Skipping route: Service address not found", "service_name", rt.Service)
			continue
		}

		p := proxy.NewSingleHostReverseProxy(cfg.ProxyConfig, svcAddr)
		finalHandler := p

		routeHandler := m.PathPrefix(rt.Prefix)
		if len(rt.Methods) > 0 {
			routeHandler = routeHandler.Methods(rt.Methods...)
		}

		routeHandler.Handler(finalHandler)

		// Log successful registration
		cfg.Log.Info(context.Background(), "Registered dynamic route",
			"prefix", rt.Prefix,
			"methods", fmt.Sprintf("%v", rt.Methods),
			"target", svcAddr)
	}

	return m
}
