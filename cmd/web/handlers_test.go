package main

import (
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
