package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// End-to-end test
func TestPing(t *testing.T) {
	//Create a new instance of our app
	app := application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}

	//Create a new test server, this starts up a HTTPS server listening on a random port
	testServer := httptest.NewTLSServer(app.routes())
	defer testServer.Close()

	//Network address that the test server is lisening on is contained in the tesetServer.URL field.
	//We can make a GET /ping request using this against the test server
	response, err := testServer.Client().Get(testServer.URL + "/ping")
	assert.NoError(t, err)

	// assert that the status code and response body
	assert.Equal(t, http.StatusOK, response.StatusCode)

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.Equal(t, "OK", string(body))

}
