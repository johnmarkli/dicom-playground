package server

import "net/http"

// HealthHandler handles the health check
type HealthHandler struct{}

// Check server health
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK"))
}
