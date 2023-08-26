package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {

	//add a command line flag for the network address
	addr := flag.String("addr", ":4000", "HTTP network address")
	//Parse reads in the CL flag value and assigns it to the addr variable. Must be called before using the addr variable.
	flag.Parse()

	//create a logger for writing information messages, which will write to std out
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	//create a logger for writing error messages, which will write to std err
	//log.Lshortfile flag will include the file name and the line number
	errLog:= log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

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

	infoLog.Printf("Starting Server on %s", *addr)

	//start a new web server
	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		errLog.Fatal(err)
	}
}
