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

	//Comman line Flag for netowrok address
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	//Custom Logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog:= log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Initialise a new instance of application contiaining the dependencies
	app := &application{
		infoLog:  infoLog,
		errorLog: errLog,
	}

	//Routes
	mux := app.routes()

	//Intialise a new Http server and set the address, handler and errorLog fields so that these are used instead of Go's default HTTP server settings
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
