package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/stasd82/la21-chars/domain/web/metrics"
	"github.com/stasd82/tux"
)

func Panics() (c tux.Circuit) {
	c = func(in tux.Route) (out tux.Route) {
		out = func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))

					metrics.AddPanics(ctx)
				}
			}()

			return in(ctx, w, r)
		}
		return
	}
	return
}
