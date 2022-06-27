package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/stasd82/tux"
	"go.uber.org/zap"
)

type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

func APIMux(cfg APIMuxConfig) *tux.Tux {

	t := tux.New(cfg.Shutdown)

	return t
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

	return mux
}
