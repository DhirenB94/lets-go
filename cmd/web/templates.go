package main

import (
	"html/template"
	"path/filepath"
	"time"

	models "dhiren.brahmbhatt/snippetbox/pkg"
	"dhiren.brahmbhatt/snippetbox/pkg/forms"
)

// Pass the new custom Form struct
type templateData struct {
	AuthenticatedUser *models.User
	CSRFToken         string
	CurrentYear       int
	Flash             string
	Forms             *forms.Form
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
}

// humanDate returns a formatted easily readable string of the date
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// Initialise a template.FuncMap object and store it in a global variable.
// This is a string-keyed map which acts as a lookup between the names of of custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	//Initialise a new map to act as the cache
	cache := map[string]*template.Template{}

	//get a slice of all the files with the page.tmpl extenstion
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		//extract the file name (e.g home.page.tmpl) from the full path
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set before we call the ParseFiles() method.
		// This means we have to use template.New to  create an empty template set,
		// Then use the Funcs() method to register the template.FuncMap,
		// Then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		//add layout and partial templates to the template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		//add the template set to the cache
		cache[name] = ts
	}
	return cache, nil

}
