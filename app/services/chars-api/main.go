package main

import (
	"errors"
	"expvar"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
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

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	return nil
}
