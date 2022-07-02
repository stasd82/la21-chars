package mid

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stasd82/tux"
)

func Panics() (c tux.Circuit) {
	c = func(in tux.Route) (out tux.Route) {
		out = func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			defer func() {
				if rec := recover(); rec != nil {
					err = fmt.Errorf("PANIC [%v]", rec)
				}
			}()

			return in(ctx, w, r)
		}
		return
	}
	return
}
