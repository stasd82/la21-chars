package mid

import (
	"context"
	"net/http"

	"github.com/stasd82/la21-chars/domain/web/metrics"
	"github.com/stasd82/tux"
)

func Metrics() (c tux.Circuit) {
	c = func(in tux.Route) (out tux.Route) {
		out = func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			ctx = metrics.Set(ctx)
			{
				err = in(ctx, w, r)
			}
			metrics.AddRequests(ctx)
			metrics.AddGoroutines(ctx)

			if err != nil {
				metrics.AddErrors(ctx)
			}
			return
		}
		return
	}
	return
}
