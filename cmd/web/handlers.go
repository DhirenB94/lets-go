package main

import (
	"fmt"

	//"html/template"
	"net/http"
	"strconv"

	models "dhiren.brahmbhatt/snippetbox/pkg"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippetsDb.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{Snippets: snippets})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	idFromUrl := r.URL.Query().Get(":id")
	if idFromUrl == "" {
		app.notFound(w)
		return
	}
	id, err := strconv.Atoi(idFromUrl)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippetsDb.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{Snippet: snippet})

	app.infoLog.Printf("Displaying a specific snippet with ID %d...", id)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm() which adds any data in POST/PUT/PATCH request bodies to the r.PostForm map. 
	err := r.ParseForm()
	if err != nil {
        app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the relevant data fields from the r.PostForm map
	title :=  r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	id, err := app.snippetsDb.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	//Redirect to show the relevant page
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
	w.Write([]byte("create a new snippet"))
}
