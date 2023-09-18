package main

import (
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"dhiren.brahmbhatt/snippetbox/pkg/models/mock"
	"github.com/golangcollege/sessions"
	"github.com/stretchr/testify/assert"
)

// Define a custom testServer struct which anonymously embed a http.server instance
type testServer struct {
	server *httptest.Server
}

// newTestServer helper will initialise a new instnace of the custom testServer
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	//initialise a new cookies jar and add it to the client
	jar, err := cookiejar.New(nil)
	assert.NoError(t, err)
	ts.Client().Jar = jar

	//disable reditrect following for the client.
	//this func will be called after a 3xx response is recieved by the client,
	//returning the http.ErrUseLastResponse error forces it to immediately return the received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{server: ts}
}

// get method on our custom testServer type will make a GET request on the URL of the server plus desired URL and return the statusCode, headers and body
func (cs *testServer) get(t *testing.T, url string) (int, http.Header, []byte) {
	response, err := cs.server.Client().Get(cs.server.URL + url)
	assert.NoError(t, err)

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	return response.StatusCode, response.Header, body
}

func (cs *testServer) postForm(t *testing.T, url string, formData url.Values) (int, http.Header, []byte) {
	response, err := cs.server.Client().PostForm(cs.server.URL+url, formData)
	assert.NoError(t, err)

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)

	return response.StatusCode, response.Header, body
}

// newTestApplication helper returns an instance of our application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *application {

	// Create a session manager instance, with the same settings as production.
	session := sessions.New([]byte("3dSm5MnygFHh7XidAtbskXrjbwfoJcbJ"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	// Create an instance of the template cache.
	templateCache, err := newTemplateCache("./../../ui/html/")
	if err != nil {
		t.Fatal(err)
	}

	return &application{
		errorLog:      log.New(io.Discard, "", 0),
		infoLog:       log.New(io.Discard, "", 0),
		session:       session,
		snippetsDb:    &mock.MockSnippetModel{},
		userDB:        &mock.MockUserModel{},
		templateCache: templateCache,
	}
}

// Define a regular expression which captures the CSRF token value from the HTML for our user signup page.
var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='([^']+)'>`)

func extractCSRFToken(t *testing.T, body []byte) string {
	// Use the FindSubmatch method to extract the token from the HTML body.
	// Note that this returns an array with the entire matched pattern in the 1st position
	// and the values of any captured data in the subsequent positions.
	matches := csrfTokenRX.FindSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}
	return html.UnescapeString(string(matches[1]))
}
