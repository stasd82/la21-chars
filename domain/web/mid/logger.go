package mid

import (
	"context"
	"net/http"

	"github.com/stasd82/tux"
	"go.uber.org/zap"
)

func Logger(sl *zap.SugaredLogger) (c tux.Circuit) {
	c = func(in tux.Route) (out tux.Route) {

		out = func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			sl.Infow("request started")
			{
				err = in(ctx, w, r)
			}
			sl.Infow("request completed")

			return
		}
		return
	}
	return
}
