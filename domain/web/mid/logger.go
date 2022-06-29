package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/stasd82/tux"
	"go.uber.org/zap"
)

func Logger(sl *zap.SugaredLogger) (c tux.Circuit) {
	c = func(in tux.Route) (out tux.Route) {

		out = func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			v, err := tux.GetValues(ctx)
			if err != nil {
				return
			}

			sl.Infow("request started", "traceid", v.TraceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr)
			{
				err = in(ctx, w, r)
			}
			sl.Infow("request completed", "traceid", v.TraceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr, "statuscode", v.StatusCode, "since", time.Since(v.Now))

			return
		}
		return
	}
	return
}
