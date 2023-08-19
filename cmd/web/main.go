package main

import (
	"log"
	"net/http"
)

func main() {
	//initalise a new router/servemux which will map a url to a handler
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Create a file server which serves files out of the "./ui/static" directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//the prefix "/static" is stripped, so the remaining part of the URL eg "/image.jpg" 
	//is used to find the corresponding file in the "./ui/static" directory. 
	//register the file server as the handler
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Starting Server on :4000")

	//start a new web server
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
