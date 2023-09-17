package main

import (
	"context"
	"fmt"
	"net/http"

	models "dhiren.brahmbhatt/snippetbox/pkg"
	"github.com/justinas/nosurf"
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

func (app *application) requireAuthenticatedUser(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//If the user is not authenticated, redirect to login page and return to prevent subsequent handlers from being executed
		if app.autheticatedUser(r) == 0 {
			http.Redirect(w, r, "/user/login", 302)
			return
		}
		nextHandler.ServeHTTP(w, r)
	})
}

// noSurf middleware function which uses a customised CSRF cookie with some flags set
func noSurf(nextHandler http.Handler) http.Handler {
	csrfHandler := nosurf.New(nextHandler)

	csrfHandler.SetBaseCookie(http.Cookie{
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})
	return csrfHandler
}

// authenticate will fecth details for the current user from the db, based on "userID" from their session
// It will then add these details to the request context
func (app *application) authenticate(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if a userID value exists in the session.
		//If this isn't present (not athenticated)then we want to pass the original and unchanged request onto the next handler in the chain as normal
		exists := app.session.Exists(r, "userID")
		if !exists {
			nextHandler.ServeHTTP(w, r)
			return
		}

		userID := app.session.GetInt(r, "userID")
		//Get the matching userID (from the request) form the db.
		//If no record found (invalid), remove the userID from the session and pass the original unchanged request onto the next handler
		user, err := app.userDB.Get(userID)
		if err == models.ErrNoRecord {
			app.session.Remove(r, "userID")
			nextHandler.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		// At this point we know user is authenticated and valid
		//Create a copy of the original request with the users details stored in request context
		// We then pass this request copy onto the next handler
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		nextHandler.ServeHTTP(w, r.WithContext(ctx))
	})
}
