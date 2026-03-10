package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iancenry/snippetbox/internal/assert"
)

func TestPing(t *testing.T){
	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ping", nil)

	ping(rr, r)

	rs := rr.Result()
	defer rs.Body.Close()

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(body), "OK")
}