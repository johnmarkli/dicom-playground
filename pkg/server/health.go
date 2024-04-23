package server

import "net/http"

// HealthHandler handles the health check
type HealthHandler struct{}

// Check server health
//
//	@Summary		Check server health
//	@Description	Check the health of the server
//	@Tags			health
//	@Produce		plain
//	@Success		200
//	@Router			/health [get]
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK"))
}
