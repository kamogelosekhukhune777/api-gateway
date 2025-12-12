package main

import (
	"context"
	"os"
	"runtime"

	"github.com/kamogelosekhukhune777/api-gateway/internal/observability/logger"
)

func main() {
	var log *logger.Logger

	traceIDFn := func(ctx context.Context) string {
		return ""
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
	return nil
}
