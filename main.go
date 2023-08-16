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
		//let the user know what methods they are allowed in the response header
		//all header changes need to be done before the write methods are called
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
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
