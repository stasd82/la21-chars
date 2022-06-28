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
	ct       []Circuit
}

// New creates a Tux value that handle a set of routes for the application.
func New(shutdown chan os.Signal, ct ...Circuit) *Tux {
	return &Tux{
		CompatRouter: bunrouter.New().Compat(),
		shutdown:     shutdown,
		ct:           ct,
	}
}

// Finish is used to gracefully shutdown our race when an integrity issue is identified.
func (t *Tux) Finish() {
	t.shutdown <- syscall.SIGTERM
}

// AddRoute registers the route (handler) function for the given pattern.
func (t *Tux) AddRoute(verb string, group string, path string, route Route, ct ...Circuit) {

	// Wrap route specific middleware around this handler first.
	route = warmupCircuit(ct, route)

	// Add app general middleware to the handler chain.
	route = warmupCircuit(t.ct, route)

	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.TODO()

		// Call the wrapped handler function.
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
