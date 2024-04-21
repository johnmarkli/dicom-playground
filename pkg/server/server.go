package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

// Server manages the lifecycle of the dime
// Server manages the lifecycle of the dime
// server
type Server struct {
	server *http.Server
}

// New creates a new Server instance
func New() (*Server, error) {
	router := mux.NewRouter()

	hh := healthHandler{}
	router.HandleFunc("/health", hh.Check).Methods("GET")

	dh := dicomHandler{}
	dicomsRouter := router.PathPrefix("/dicoms").Subrouter()
	dicomsRouter.HandleFunc("", dh.UploadDICOM).Methods("POST")
	dicomsRouter.HandleFunc("", dh.ListDICOMs).Methods("GET")

	s := &Server{
		server: &http.Server{
			Addr:    ":8080",
			Handler: router,
		},
	}
	return s, nil
}

// Run the dime server
// Run the dime server
func (s *Server) Run() {
	slog.Info("Starting dime server")
	slog.Info("Starting dime server")
	go s.server.ListenAndServe()
}

// Shutdown the dime server
// Shutdown the dime server
func (s *Server) Shutdown() error {
	slog.Info("Shutting down dime server")
	slog.Info("Shutting down dime server")
	return s.server.Shutdown(context.Background())
}
