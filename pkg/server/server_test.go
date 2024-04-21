package server_test

import (
	"testing"

	"github.com/johnmarkli/dime/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	s, err := server.New()
	assert.NoError(t, err)
	assert.NotNil(t, s)
}
