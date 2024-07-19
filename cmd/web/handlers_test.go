package main

import (
	"net/http"
	"testing"

	"github.com/gustavodiasag/notebox/internal/assert"
)

func TestHealthCheck(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/health_check")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}
