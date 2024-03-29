package main

import (
	"Raft/src/cmd"
	"context"
	"os"
	"os/signal"
	"runtime"

	"github.com/sirupsen/logrus"
)

var onlyOneSignalHandler = make(chan struct{})

func main() {
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	ctx := setupSignalHandler()
	if err := cmd.StartServer(ctx); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

}

// setupSignalHandler returns the context which take care of any os interrupt.
func setupSignalHandler() context.Context {
	close(onlyOneSignalHandler) // panics when called twice

	var shutdownSignals = []os.Signal{os.Interrupt}
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return ctx
}
