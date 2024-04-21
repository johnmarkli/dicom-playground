package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/johnmarkli/dime/pkg/store"
)

// Server manages the lifecycle of the dime
// server
type Server struct {
	server *http.Server
}

// New creates a new Server instance
func New() (*Server, error) {
	router := mux.NewRouter()

	hh := HealthHandler{}
	router.HandleFunc("/health", hh.Check).Methods("GET")

	st := store.NewFileStore("data")
	dh := NewDICOMHandler(st)
	dicomsRouter := router.PathPrefix("/dicoms").Subrouter()
	dicomsRouter.HandleFunc("", dh.Upload).Methods("POST")
	dicomsRouter.HandleFunc("", dh.List).Methods("GET")
	dicomsRouter.HandleFunc("/{id}", dh.Read).Methods("GET")
	dicomsRouter.HandleFunc("/{id}/attributes", dh.Attributes).Methods("GET")
	dicomsRouter.HandleFunc("/{id}/image", dh.Image).Methods("GET")

	s := &Server{
		server: &http.Server{
			Addr:    ":8080",
			Handler: router,
		},
	}
	return s, nil
}

// Run the dime server
func (s *Server) Run() {
	slog.Info("Starting dime server")
	go s.server.ListenAndServe()
}

// Shutdown the dime server
func (s *Server) Shutdown() error {
	slog.Info("Shutting down dime server")
	return s.server.Shutdown(context.Background())
}
