package server_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/johnmarkli/dime/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health", nil)
	hh := server.HealthHandler{}
	hh.Check(w, r)
	defer w.Result().Body.Close()

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	b, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.Equal(t, "OK", string(b))
}
