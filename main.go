// Package main is the dime server
package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/johnmarkli/dime/pkg/server"
)

// Environment:
//
//	DIME_PORT
//	    int - port for server to listen on
//	DIME_DATA_DIR
//	    string - directory to save data to the file system

//	@title			dime API
//	@version		1.0
//	@description	dime is a small web service designed to work with DICOM files.

//	@contact.name	John Li
//	@contact.url	http://www.swagger.io/support
//	@contact.email	johnmarkli@gmail.com

func main() {

	// Setup exit code for graceful shutdown
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	// Start dime server
	s, err := server.New()
	if err != nil {
		slog.Error(err.Error())
		exitCode = 1
		return
	}
	s.Run()
	defer s.Shutdown()

	// Handle signals
	quit := make(chan os.Signal, 1)
	defer close(quit)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
