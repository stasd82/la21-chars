package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/stasd82/la21-chars/foundation/logger"
	"go.uber.org/zap"
)

func main() {

	// Construct the application logger
	sl, err := logger.New("CHARS-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer sl.Sync()

	// Perform the startup and shutdown sequence
	if err := run(sl); err != nil {
		sl.Errorw("startup", "ERROR", err)
		sl.Sync()
		os.Exit(1)
	}
}

func run(sl *zap.SugaredLogger) error {

	sl.Infow("It's works!")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	return nil
}
