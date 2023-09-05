package main

import "net/http"

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
