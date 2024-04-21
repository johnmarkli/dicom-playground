package server

import "net/http"

// HealthHandler handles the health check
type HealthHandler struct{}

// Check the server's health
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
