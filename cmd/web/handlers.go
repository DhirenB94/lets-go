package main

import (
	"fmt"

	//"html/template"
	"net/http"
	"strconv"

	models "dhiren.brahmbhatt/snippetbox/pkg"
	"dhiren.brahmbhatt/snippetbox/pkg/forms"
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

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: snippet,
	})

	app.infoLog.Printf("Displaying a specific snippet with ID %d...", id)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm() which adds any data in POST/PUT/PATCH request bodies to the r.PostForm map.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//Initialise a new Form struct
	form := forms.NewForm(r.PostForm)

	//perform some basic validation, if validation fails add it to the validationErrors map
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	// If there are any errors, re-display the create snippet page with previously submitted data and the validation errors
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{
			Forms: form,
		})
		return
	}

	//Insert the now validated form data, instead of the users unvalidated input into the form
	id, err := app.snippetsDb.Insert(form.FormData.Get("title"), form.FormData.Get("content"), form.FormData.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	//Add a flash confirmation message
	//If there's no existing session for the current user (or their session has expired) then a new, empty, session will automatically be created by the session middleware.
	app.session.Put(r, "flash", "Snippet successfully created")

	//Redirect to show the relevant page
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Forms: forms.NewForm(nil),
	})

}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Forms: forms.NewForm(nil),
	})
	fmt.Fprintln(w, "Display the signup user form")
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	//Parse the form data
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//validate the user input
	form := forms.NewForm(r.PostForm)
	form.Required("name", "email", "password")
	form.MinLength("password", 10)
	form.MatchesPattern("email", forms.EmailRX)

	//if there are any errors, re-display the signup form
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{
			Forms: form,
		})
		return
	}

	//Insert the new user record into the DB.
	//If the email already exists, add an error messsage and redeisplay the form
	err = app.userDB.Insert(form.FormData.Get("name"), form.FormData.Get("email"), form.FormData.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.FormErrors.Add("email", "email already exists")
		app.render(w, r, "signup.page.tmpl", &templateData{Forms: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	//Add a flash confirmation message
	app.session.Put(r, "flash", "Signup successfull, please login")

	//Redirect to the login page
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display the login user form")
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Login user")
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout user")
}
