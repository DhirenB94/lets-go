package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

//Define struct to hold application-wide dependencies 
type application struct {
	infoLog *log.Logger
	errorLog *log.Logger
}

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

	//Initialise a new instance of application contiainig the dependencies
	app := &application{
		infoLog:  infoLog,
		errorLog: errLog,
	}

	//initalise a new router/servemux which will map a url to a handler
	mux := http.NewServeMux()

	//swap the rouitng declerations as handlers are now mehtods on application
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Create a file server which serves files out of the "./ui/static" directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//the prefix "/static" is stripped, so the remaining part of the URL eg "/image.jpg"
	//is used to find the corresponding file in the "./ui/static" directory.
	//register the file server as the handler
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	//intialise a new Http server and set the address, handler and errorLog fields so that these are used instead of Go's default HTTP server settings
	srv := http.Server{
		Addr: *addr,
		Handler: mux,
		ErrorLog: errLog,
	}

	infoLog.Printf("Starting Server on %s", *addr)

	//call the listen and serve method on our new htpp server
	err := srv.ListenAndServe()
	if err != nil {
		errLog.Fatal(err)
	}
}
