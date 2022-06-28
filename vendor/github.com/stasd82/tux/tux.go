package tux

import (
	"context"
	"net/http"
	"os"
	"syscall"

	"github.com/uptrace/bunrouter"
)

// The Route type is an adapter to allow the use of our advanced functions as HTTP handlers.
type Route func(context.Context, http.ResponseWriter, *http.Request) error

type Tux struct {
	*bunrouter.CompatRouter
	shutdown chan os.Signal
}

func New(c chan os.Signal) *Tux {
	return &Tux{
		CompatRouter: bunrouter.New().Compat(),
		shutdown:     c,
	}
}

// Finish is used to gracefully shutdown our race when an integrity issue is identified.
func (t *Tux) Finish() {
	t.shutdown <- syscall.SIGTERM
}

// AddRoute registers the route (handler) function for the given pattern.
func (t *Tux) AddRoute(verb string, group string, path string, route Route) {

	h := func(w http.ResponseWriter, r *http.Request) {

		ctx := context.TODO()
		if err := route(ctx, w, r); err != nil {
			t.Finish()
			return
		}

	}

	team := t.NewGroup("/" + group)

	switch verb {
	case http.MethodGet:
		team.GET(path, h)
	case http.MethodPost:
		team.POST(path, h)
	case http.MethodPut:
		team.PUT(path, h)
	case http.MethodDelete:
		team.DELETE(path, h)
	}
}
