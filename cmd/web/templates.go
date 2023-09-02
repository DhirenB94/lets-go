package main

import (
	"html/template"
	"path/filepath"

	models "dhiren.brahmbhatt/snippetbox/pkg"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
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

		//parse the page template file into a template set
		ts, err := template.ParseFiles(page)
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
