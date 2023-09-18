package main

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// End-to-end test
func TestPing(t *testing.T) {
	//Create a new instance of our app
	app := newTestApplication(t)

	//Create a new test server, this starts up a HTTPS server listening on a random port
	testServer := newTestServer(t, app.routes())
	defer testServer.server.Close()

	//make the get request against the testServer and get the responses
	statusCode, _, body := testServer.get(t, "/ping")

	// assert
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, "OK", string(body))

}

func TestShowSnippet(t *testing.T) {
	//Create a new instance of our app struct which uses mocked dependencies
	app := newTestApplication(t)

	//Create a new test server, this starts up a HTTPS server listening on a random port
	testServer := newTestServer(t, app.routes())
	defer testServer.server.Close()

	//Create some test scenarios
	//Our mocks.mockSnippet.Get() is returning the mockSnippet if ID is 1 and returning ErrNoRecord if anything else
	tests := []struct {
		name         string
		url          string
		expectedCode int
		expectedBody []byte
	}{
		{
			name:         "Valid ID",
			url:          "/snippet/1",
			expectedCode: http.StatusOK,
			expectedBody: []byte("mock"),
		},
		{
			name:         "Invalid ID",
			url:          "/snippet/2",
			expectedCode: http.StatusNotFound,
			expectedBody: nil,
		},
		{
			name:         "Negative ID",
			url:          "/snippet/-1",
			expectedCode: http.StatusNotFound,
			expectedBody: nil,
		},
		{
			name:         "Deciamal ID",
			url:          "/snippet/1.5",
			expectedCode: http.StatusNotFound,
			expectedBody: nil,
		},
		{
			name:         "String ID",
			url:          "/snippet/foo",
			expectedCode: http.StatusNotFound,
			expectedBody: nil,
		},
		{
			name:         "Empty ID",
			url:          "/snippet/",
			expectedCode: http.StatusNotFound,
			expectedBody: nil,
		},
		{
			name:         "Trailing Slash",
			url:          "/snippet/1/",
			expectedCode: http.StatusNotFound,
			expectedBody: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, _, body := testServer.get(t, tt.url)

			assert.Equal(t, tt.expectedCode, statusCode)

			if !bytes.Contains(body, tt.expectedBody) {
				t.Errorf("want body to contain %q", tt.expectedBody)
			}
		})
	}

}

func TestSignupUser(t *testing.T) {
	//Create a new instance of our app struct which uses mocked dependencies
	app := newTestApplication(t)

	//Create a new test server, this starts up a HTTPS server listening on a random port
	testServer := newTestServer(t, app.routes())
	defer testServer.server.Close()

	//Make a GET /user/signup request 1st
	//This will return a response which has a CSRF cookie in the header
	//And a CSRF token in the HTML response body
	getCode, _, body := testServer.get(t, "/user/signup")
	assert.Equal(t, http.StatusOK, getCode)

	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantBody     []byte
	}{
		{
			name:         "Valid Submission",
			userName:     "Bob",
			userEmail:    "bob@mail.com",
			userPassword: "validPassword",
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
			wantBody:     []byte(""),
		},
		{
			name:         "Empty name",
			userName:     "",
			userEmail:    "bob@mail.com",
			userPassword: "validPassword",
			csrfToken:    csrfToken,
			wantCode:     http.StatusOK,
			wantBody:     []byte("This field cannot be blank"),
		},
		{
			name:         "Empty email",
			userName:     "Bob",
			userEmail:    "",
			userPassword: "validPassword",
			csrfToken:    csrfToken,
			wantCode:     http.StatusOK,
			wantBody:     []byte("This field cannot be blank"),
		},
		{
			name:         "Empty password",
			userName:     "Bob",
			userEmail:    "bob@mail.com",
			userPassword: "",
			csrfToken:    csrfToken,
			wantCode:     http.StatusOK,
			wantBody:     []byte("This field cannot be blank"),
		},
		{
			name:         "inavlid email",
			userName:     "Bob",
			userEmail:    "bobmail.com",
			userPassword: "validPassword",
			csrfToken:    csrfToken,
			wantCode:     http.StatusOK,
			wantBody:     []byte("This email is invalid"),
		},
		{
			name:         "Duplicate email",
			userName:     "Bob",
			userEmail:    "dupe@example.com",
			userPassword: "validPassword",
			csrfToken:    csrfToken,
			wantCode:     http.StatusOK,
			wantBody:     []byte("email already exists"),
		},
		{
			name:         "inavlid password ",
			userName:     "Bob",
			userEmail:    "bob@mail.com",
			userPassword: "short",
			csrfToken:    csrfToken,
			wantCode:     http.StatusOK,
			wantBody:     []byte("This field is too short"),
		},
		{
			name:         "inavlid crsf token ",
			userName:     "",
			userEmail:    "bob@.com",
			userPassword: "",
			csrfToken:    "invalidToken",
			wantCode:     http.StatusBadRequest,
			wantBody:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formData := url.Values{}
			formData.Add("name", tt.userName)
			formData.Add("email", tt.userEmail)
			formData.Add("password", tt.userPassword)
			formData.Add("csrf_token", tt.csrfToken)

			codePost, _, body := testServer.postForm(t, "/user/signup", formData)

			assert.Equal(t, tt.wantCode, codePost)
			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}
