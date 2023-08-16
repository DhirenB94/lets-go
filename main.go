package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	 //ensure the URL path exactly matches "/" so home handler is only called on the homepage
    if r.URL.Path != "/" {
        http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from snippetbox"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		//lightweight helper function which takes a given message and status code, then calls the w.WriteHeader() and w.Write() methods behind-the-scenes.
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}

func main() {
	//initalise a new router/servemux which will map a url to a handler
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	log.Println("Starting Server on :4000")

	//start a new web server
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}

}
