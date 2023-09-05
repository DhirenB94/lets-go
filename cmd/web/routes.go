package main

import "net/http"

// Because we want this middleware to act on every request that is recieved
// We need to execute it before a request hits the servermux
// So, we need to wrap secureHeaders middleware around servemux

// Update the signature of routes() so that it returns a handler instead of a servemux
func (app *application) routes() http.Handler {
	//Initalise a new router/servemux
	mux := http.NewServeMux()

	//Routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	//File Server & register it as a handler
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return secureHeaders(mux)
}
