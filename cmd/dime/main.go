// Package main ...
package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/johnmarkli/dime/pkg/server"
)

func main() {

	// setup exit code for graceful shutdown
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	s, err := server.New()
	if err != nil {
		slog.Error(err.Error())
		exitCode = 1
		return
	}
	s.Run()
	defer s.Shutdown()

	quit := make(chan os.Signal, 1)
	defer close(quit)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
