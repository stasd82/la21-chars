package handlers

import (
	"context"
	"errors"
	"expvar"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/stasd82/la21-chars/app/services/chars-api/handlers/debug/checkgrp"
	v1 "github.com/stasd82/la21-chars/domain/web/v1"
	"github.com/stasd82/la21-chars/domain/web/v1/mid"
	"github.com/stasd82/tux"
	"go.uber.org/zap"
)

type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

func APIMux(cfg APIMuxConfig) http.Handler {
	var mux *tux.Tux

	if mux == nil {
		mux = tux.New(
			cfg.Shutdown,
			mid.Logger(cfg.Log),
			mid.Errors(cfg.Log),
		)
	}

	// Load the routes for different versions of the API.
	bindV1(mux, cfg)

	return mux
}

func bindV1(t *tux.Tux, cfg APIMuxConfig) {
	const version = "v1"

	// Test handler for the development and testing.
	route := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		if rand.Intn(100)%2 == 0 {
			// return errors.New("untrusted error")
			// return tux.NewShutdownError("going down")
			return v1.NewRequestErr(errors.New("trusted error"), http.StatusBadGateway)
		}

		msg := struct {
			Message string
		}{
			Message: "yey",
		}
		return tux.Respond(ctx, w, msg, http.StatusOK)
	}

	t.AddRoute(http.MethodGet, version, "/test", route)
}

func debugStdLibMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

func DebugMux(build string, log *zap.SugaredLogger) http.Handler {
	mux := debugStdLibMux()

	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}
