package handlers

import (
	"context"
	"encoding/json"
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/stasd82/la21-chars/app/services/chars-api/handlers/debug/checkgrp"
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
		)
	}

	// Load the routes for different versions of the API.
	bindV1(mux, cfg)

	return mux
}

func bindV1(tux *tux.Tux, cfg APIMuxConfig) {
	const version = "v1"

	// Test handler for the development and testing.
	route := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		st := struct {
			Message string
		}{
			Message: "hello world",
		}
		return json.NewEncoder(w).Encode(st)
	}

	tux.AddRoute(http.MethodGet, version, "/test", route)
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
