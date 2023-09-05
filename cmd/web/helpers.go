package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// serverError writes an error to the stack trace to the customLogger, then sends a generic 500 internal server error to the user
func (app *application) severError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends a specific status code and corresponding desription to the user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound is a wrapper of the clientErorr to send a 404 not found response to the user
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

//addDefaultYear
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData  {
	if td == nil {
        td = &templateData{}
    }
	td.CurrentYear = time.Now().Year()
	return td
}

//render looks up the template set in the cache
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	//Retrieve the appropriate template set based on the name of the page
	ts, ok := app.templateCache[name]
	if !ok {
		app.severError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
    //Initialise a new buffer
    buf := new(bytes.Buffer)

	//add the default data
	app.addDefaultData(td, r)

	//execute the template set with the current year injected
	err := ts.Execute(buf, td)
	if err != nil {
		app.severError(w, err)
		return
	}

	//Write the contents of the buffer into the resposnse writer
	buf.WriteTo(w)
}
