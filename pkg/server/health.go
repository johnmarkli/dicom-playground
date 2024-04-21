package server

import "net/http"

type healthHandler struct {
	server *Server
}

func (h *healthHandler) Check(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
