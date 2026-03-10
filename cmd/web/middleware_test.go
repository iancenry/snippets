package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iancenry/snippetbox/internal/assert"
)

func TestSecureHeader(t *testing.T){
	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(rr, r)

	rs:= rr.Result()
	defer rs.Body.Close()


	expectedValue := "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)

	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)

	expectedValue = "default-src 'self'; style-src 'self' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	// Check that the next handler in the chain is called and writes the expected response body
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	
	assert.Equal(t, string(body), "OK")
}