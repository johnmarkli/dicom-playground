// Package server defines the server for dime
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/johnmarkli/dime/pkg/store"
)

const (
	defaultPort    = 8080
	defaultDataDir = "data"
)

// Server manages the lifecycle of the dime server
type Server struct {
	server *http.Server
}

// New creates a new Server instance
//
// Environment:
//
//	DIME_PORT
//	    int - port for server to listen on
//	DIME_DATA_DIR
//	    string - directory to save data to the file system
func New() (*Server, error) {
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	hh := HealthHandler{}
	router.HandleFunc("/health", hh.Check).Methods("GET")

	st, err := store.NewFileStore(getDataDir())
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}
	dh := NewDICOMHandler(st)
	dicomsRouter := router.PathPrefix("/dicoms").Subrouter()
	dicomsRouter.HandleFunc("", dh.Upload).Methods("POST")
	dicomsRouter.HandleFunc("", dh.List).Methods("GET")
	dicomsRouter.HandleFunc("/{id}", dh.Read).Methods("GET")
	dicomsRouter.HandleFunc("/{id}/attributes", dh.Attributes).Methods("GET")
	dicomsRouter.HandleFunc("/{id}/image", dh.Image).Methods("GET")

	s := &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", getPort()),
			Handler: router,
		},
	}
	return s, nil
}

// Run the dime server
func (s *Server) Run() {
	slog.Info("Starting dime server", slog.String("on", s.server.Addr))
	go s.server.ListenAndServe()
}

// Shutdown the dime server
func (s *Server) Shutdown() error {
	slog.Info("Shutting down dime server")
	return s.server.Shutdown(context.Background())
}

// Server returns the http server
func (s *Server) Server() *http.Server {
	return s.server
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request received",
			slog.String("method", r.Method),
			slog.String("request", r.RequestURI))
		next.ServeHTTP(w, r)
	})
}

func getPort() int {
	port := defaultPort
	if val, ok := os.LookupEnv("DIME_PORT"); ok {
		if p, err := strconv.Atoi(val); err == nil {
			port = p
		}
	}
	return port
}

func getDataDir() string {
	dir := defaultDataDir
	if val, ok := os.LookupEnv("DIME_DATA_DIR"); ok {
		dir = val
	}
	return dir
}
