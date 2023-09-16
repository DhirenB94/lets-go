package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	//Initalise a new router/servemux
	mux := pat.New()

	//create a standard middleware chain
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	//For the sessions to work, we need to wrap our application routes with the middleware provided by the Session.Enable() method. 
	//This middleware loads and saves session data to and from the session cookie with every HTTP request and response as appropriate.
	//We only need this middleware to act on our dynamic routes, the static route onlky serves static files and so does not need stateful behaviour
	//create a dynamic middleware chain
	dynamicMiddleware := alice.New(app.session.Enable)

	//Routes
	//Uodate the routes using the dynamicmiddlewarechain
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	//File Server & register it as a handler
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	//User related routes
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))


	return standardMiddleware.Then(mux)
}
