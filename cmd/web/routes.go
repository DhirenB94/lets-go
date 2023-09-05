package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	//Initalise a new router/servemux
	mux := http.NewServeMux()

	//create a standard middleware chain
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	//Routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	//File Server & register it as a handler
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
