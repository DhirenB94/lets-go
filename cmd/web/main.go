package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"dhiren.brahmbhatt/snippetbox/pkg/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// Define struct to hold application-wide dependencies
type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	snippetsDb *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	//Command line Flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	//Defenining a new CL flag for MySQL Data Source Name string
	dsn := flag.String("dsn", "web:dhiren@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	//Custom Logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Initialise a connection pool to the database
	db, err := sql.Open("mysql", *dsn)
	if err != nil {
		errLog.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		errLog.Fatal(err)
	}
	defer db.Close()

	//Initialise a new template cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errLog.Fatal(err)
	}

	//Initialise a new instance of application contiaining the dependencies
	app := &application{
		infoLog:  infoLog,
		errorLog: errLog,
		snippetsDb: &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	//Routes
	mux := app.routes()

	//Intialise a new Http server and set the address, handler and errorLog fields so that these are used instead of Go's default HTTP server settings
	srv := http.Server{
		Addr:     *addr,
		Handler:  mux,
		ErrorLog: errLog,
	}

	infoLog.Printf("Starting Server on %s", *addr)

	//call the listen and serve method on our new htpp server
	err = srv.ListenAndServe()
	if err != nil {
		errLog.Fatal(err)
	}
}
