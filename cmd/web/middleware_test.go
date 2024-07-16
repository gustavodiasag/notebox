package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gustavodiasag/notebox/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Mock HTTP handler to be passed into the middleware.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	expected := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expected)

	expected = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expected)

	expected = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expected)

	expected = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expected)

	expected = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expected)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
