package main

import (
	"context"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/kamogelosekhukhune777/api-gateway/internal/config"
	"github.com/kamogelosekhukhune777/api-gateway/internal/observability/logger"
	"github.com/kamogelosekhukhune777/api-gateway/internal/proxy"
	"github.com/kamogelosekhukhune777/api-gateway/internal/router"
	"github.com/kamogelosekhukhune777/api-gateway/internal/server"
	"github.com/kamogelosekhukhune777/api-gateway/internal/trace"
)

func main() {
	var log *logger.Logger

	traceIDFn := func(ctx context.Context) string {
		return trace.GetTraceID(ctx).String()
	}

	log = logger.New(os.Stdout, logger.LevelInfo, "API-GATEWAY", traceIDFn)

	// -------------------------------------------------------------------------

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	cfg, err := config.LoadConfig("./configs/gateway.yaml")
	if err != nil {
		return err
	}

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	log.BuildInfo(ctx)

	expvar.NewString("build").Set(cfg.Build)

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	pcfg := proxy.Config{
		Log:                   log,
		DialTimeout:           cfg.DialTimeout,
		ResponseHeaderTimeout: cfg.ResponseHeaderTimeout,
		KeepAlive:             cfg.KeepAlive,
		MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
	}

	rcfg := router.Config{
		Log:         log,
		Services:    cfg.Services,
		Routes:      cfg.RouterConfig.Routes,
		ProxyConfig: pcfg,
	}

	scfg := server.Config{
		RouterConfig: rcfg,
	}

	webAPI := server.NewServer(&scfg)

	api := http.Server{
		Addr:         cfg.ServerConfig.APIHost,
		Handler:      webAPI,
		ReadTimeout:  cfg.ServerConfig.ReadTimeout,
		WriteTimeout: cfg.ServerConfig.WriteTimeout,
		IdleTimeout:  cfg.ServerConfig.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.ServerConfig.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
