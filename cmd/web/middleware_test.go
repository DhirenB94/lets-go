package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecureHeaders(t *testing.T) {
	//initialise a dummy http.Request and httptetest.ResponseRecorder
	request, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()

	//initalise a mock http handler that we can pass to our secureheaders middleware, which writes back a 200 and "OK in the response body"
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	//pass the mockhandler to the middleware we want to test
	secureHeaders(mockHandler).ServeHTTP(response, request)

	//get the results
	result := response.Result()

	//assert that the middleware has correctly set the headers on the response
	assert.Equal(t, "deny", result.Header.Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", result.Header.Get("X-XSS-Protection"))

	//check that the middleware has correctly called the next handler in line and the response and body are as expected
	assert.Equal(t, http.StatusOK, result.StatusCode)
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	assert.NoError(t, err)
	assert.Equal(t, string(body), "OK")
}
