package main

import (
	"net/http"
	"net/url"
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

func TestUserSignup(t *testing.T) {
	app := newTestApplication()

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, string(body))

	const (
		name = "Test User"
		email = "testuser@example.com"
		password = "correctpassword"
		formTag = `<form action='/user/signup' method='POST' novalidate>`
	)

	tests := []struct {	
		name string
		userName string
		userEmail string
		userPassword string
		csrfToken string
		expectedCode int
		expectedFormTag string
	}{
		{
			name: "valid form submission",
			userName: name,
			userEmail: email,
			userPassword: password,
			csrfToken: csrfToken,
			expectedCode: http.StatusSeeOther,
		},
		{
			name: "Invalid CSRF token",
			userName: name,
			userEmail: email,
			userPassword: password,
			csrfToken: "invalidtoken",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Empty name field",
			userName: "",
			userEmail: email,
			userPassword: password,
			csrfToken: csrfToken,
			expectedCode: http.StatusUnprocessableEntity,
			expectedFormTag: formTag,
		},
		{
			name: "Empty email field",
			userName: name,
			userEmail: "",
			userPassword: password,
			csrfToken: csrfToken,
			expectedCode: http.StatusUnprocessableEntity,
			expectedFormTag: formTag,
		},
		{
			name: "Empty password field",
			userName: name,
			userEmail: email,
			userPassword: "",
			csrfToken: csrfToken,
			expectedCode: http.StatusUnprocessableEntity,
			expectedFormTag: formTag,
		},
		{
			name: "Invalid email format",
			userName: name,
			userEmail: "invalid-email",
			userPassword: password,
			csrfToken: csrfToken,
			expectedCode: http.StatusUnprocessableEntity,
			expectedFormTag: formTag,
		},
		{
			name: "Short password",
			userName: name,
			userEmail: email,
			userPassword: "short",
			csrfToken: csrfToken,
			expectedCode: http.StatusUnprocessableEntity,
			expectedFormTag: formTag,
		},
		{
			name: "Duplicate email",
			userName: name,
			userEmail: "dupe@example.com",
			userPassword: password,
			csrfToken: csrfToken,
			expectedCode: http.StatusUnprocessableEntity,
			expectedFormTag: formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formData := url.Values{
				"name": {tt.userName},
				"email": {tt.userEmail},
				"password": {tt.userPassword},
				"csrf_token": {tt.csrfToken},
			}

			code, _, body := ts.postForm(t, "/user/signup", formData)
			assert.Equal(t, code, tt.expectedCode)

			if tt.expectedFormTag != "" {
				assert.StringContains(t, string(body), tt.expectedFormTag)
			}
		})
	}
		
}