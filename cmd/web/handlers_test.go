package main

import (
	"net/http"
	"testing"

	"github.com/iancenry/snippetbox/internal/assert"
)

func TestPing(t *testing.T){
	app := newTestApplication()

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")
	
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, string(body), "OK")
}


func TestSnippetView(t *testing.T) {
	app := newTestApplication()

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
	name string
	urlPath string
	expectedCode int
	expectedBody string
}{
	{
		name: "valid UUID",
		urlPath: "/snippet/view/00000000-0000-0000-0000-000000000001",
		expectedCode: http.StatusOK,
		expectedBody: "This is the first test snippet.",
	},
	{
		name: "non-existent UUID",
		urlPath: "/snippet/view/00000000-0000-0000-0000-000000000999",
		expectedCode: http.StatusNotFound,
	},
	{
		name: "invalid UUID format",
		urlPath: "/snippet/view/invalid-uuid",
		expectedCode: http.StatusNotFound,
	},
	{
		name: "empty UUID",
		urlPath: "/snippet/view/",
		expectedCode: http.StatusNotFound,
	},
	{
		name: "UUID with wrong format",
		urlPath: "/snippet/view/12345",
		expectedCode: http.StatusNotFound,
	},
	{
		name: "String ID instead of UUID",
		urlPath: "/snippet/view/hello",
		expectedCode: http.StatusNotFound,
	},
}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)
			assert.Equal(t, code, tt.expectedCode)
			assert.StringContains(t, string(body), tt.expectedBody)
		})
	}
}