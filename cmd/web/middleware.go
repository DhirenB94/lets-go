package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		nextHandler.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s-%s, %s, %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		nextHandler.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//create a deffered function which will always be run in the event of a panic
		defer func() {
			//use built in recover function to check if there has been a panic or not
			if err := recover(); err != nil {
				//Set a Connection close header on the response
				w.Header().Set("Connection", "close")
				//call our helper method to return a 500 internal server error response
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		nextHandler.ServeHTTP(w, r)
	})
}
