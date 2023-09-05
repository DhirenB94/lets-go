package main

import "net/http"

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

	//logRequest <-> secureHeaders <-> mux <-> application handlers
	return app.logRequest(secureHeaders(mux))
}
