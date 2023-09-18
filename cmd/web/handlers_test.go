package main

import (
	"bytes"
	"net/http"
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
			expectedBody: []byte(""),
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
