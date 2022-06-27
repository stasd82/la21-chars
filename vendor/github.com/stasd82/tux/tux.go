package tux

import (
	"os"
	"syscall"

	"github.com/uptrace/bunrouter"
)

type Tux struct {
	*bunrouter.Router
	shutdown chan os.Signal
}

func New(c chan os.Signal) *Tux {
	return &Tux{
		Router:   bunrouter.New(),
		shutdown: c,
	}
}

func (t *Tux) Finish() {
	t.shutdown <- syscall.SIGTERM
}
