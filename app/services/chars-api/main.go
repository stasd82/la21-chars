package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/stasd82/la21-chars/app/services/chars-api/handlers"
	"github.com/stasd82/la21-chars/foundation/logger"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

var build = "develop"

func main() {

	// Construct the application logger
	sl, err := logger.New("CHARS-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer sl.Sync()

	// Perform the startup and shutdown sequence
	if err := run(sl); err != nil {
		sl.Errorw("startup", "ERROR", err)
		sl.Sync()
		os.Exit(1)
	}
}

func run(sl *zap.SugaredLogger) error {

	// -----------------------------------------------------------
	// GOMAXPROCS

	options := maxprocs.Logger(sl.Infof)

	if _, err := maxprocs.Set(options); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	sl.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -----------------------------------------------------------
	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Web service in Go using Kubernetes",
		},
	}

	const prefix = "CHARS"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// -----------------------------------------------------------
	// App Starting

	sl.Infow("starting service", "version", build)
	defer sl.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	sl.Infow("startup", "config", out)

	expvar.NewString("build").Set(build)

	// -----------------------------------------------------------
	// Start Debug Service

	sl.Infow("startup", "status", "debug router started", "host", cfg.Web.DebugHost)

	debugMux := handlers.DebugStandardLibraryMux()

	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debugMux); err != nil {
			sl.Errorw("shutdown", "status", "debug router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	// -----------------------------------------------------------
	// Start API Service

	sl.Infow("startup", "status", "init API support")

	svcErrors := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      nil,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(sl.Desugar()),
	}

	go func() {
		sl.Infow("startup", "status", "API router started", "host", cfg.Web.APIHost)
		svcErrors <- api.ListenAndServe()
	}()

	// -----------------------------------------------------------
	// Graceful Shutdown

	select {
	case err := <-svcErrors:
		return fmt.Errorf("svc error: %w", err)
	case sig := <-shutdown:
		sl.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer sl.Infow("shutdown", "status", "shutdown completed", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
