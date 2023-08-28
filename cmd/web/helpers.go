package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)
//serverError writes an error to the stack trace to the customLogger, then sends a generic 500 internal server error to the user
func (app *application) severError (w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

//clientError sends a specific status code and corresponding desription to the user
func (app *application) clientError(w http.ResponseWriter, status int)  {
	http.Error(w, http.StatusText(status), status)
}

//notFound is a wrapper of the clientErorr to send a 404 not found response to the user
func (app *application) notFound(w http.ResponseWriter)  {
	app.clientError(w, http.StatusNotFound)
}