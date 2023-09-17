package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Define a custom testServer struct which anonymously embed a http.server instance
type testServer struct {
	server *httptest.Server
}

// newTestServer helper will initialise a new instnace of the custom testServer
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	return &testServer{server: ts}
}

// get method on our custom testServer type will make a GET /ping request on the test server and return the statusCode, headers and body
func (cs *testServer) get(t *testing.T, url string) (int, http.Header, []byte) {
	response, err := cs.server.Client().Get(cs.server.URL + url)
	assert.NoError(t, err)

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	return response.StatusCode, response.Header, body
}

// newTestApplication helper returns an instance of our application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *application {
	return &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}
}
