package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
        "./ui/html/home.page.tmpl",
        "./ui/html/base.layout.tmpl",
    }
	//read the template file into the template set
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "unable to parse", http.StatusInternalServerError)
		return
	}
	//execute  to write the template content as the response body
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
        http.Error(w, "unable to execute template", 500)
		return
	}
	w.Write([]byte("Hello from snippetbox"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	idFromUrl := r.URL.Query().Get("id")
	if idFromUrl == "" {
		http.Error(w, "no id entered", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idFromUrl)
	if err != nil || id < 1 {
		http.Error(w, "invalid id", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
