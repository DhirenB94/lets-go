package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	models "dhiren.brahmbhatt/snippetbox/pkg/models"
	"dhiren.brahmbhatt/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type snippetDB interface {
	Insert(title, content, expires string) (int, error)
	Get(id int) (*models.Snippet, error)
	Latest() ([]*models.Snippet, error)
}

type userDB interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Get(id int) (*models.User, error)
}

// Define struct to hold application-wide dependencies
type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	session       *sessions.Session
	snippetsDb    snippetDB
	userDB        userDB
	templateCache map[string]*template.Template
}

// Define contextKey type and contextKeyUser variable so that we have a unique key that we can use to store and retrieve the user details from the request context
type contextKey string

var contextKeyUser = contextKey("user")

func main() {

	//Command line Flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	//MySQL Data Source Name
	dsn := flag.String("dsn", "web:dhiren@/snippetbox?parseTime=true", "MySQL data source name")
	//Define a new CL flag for the session secret (a random key used to encrypt and authenticate session cookies, 32 bytes long)
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Session Secret")
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

	//Initialise a new session manager, pass in the secret key
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true //Set the Secure flag on our session cookies

	//Initialise a new instance of application contiaining the dependencies
	app := &application{
		infoLog:       infoLog,
		errorLog:      errLog,
		session:       session,
		snippetsDb:    &mysql.SnippetModel{DB: db},
		userDB:        &mysql.UserModel{DB: db},
		templateCache: templateCache,
	}

	// Initialise a tls.Config struct
	tlsConfig := tls.Config{
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		PreferServerCipherSuites: true,
	}

	//Routes
	mux := app.routes()

	//Intialise a new Http server and set the address, handler, errorLog and TLS fields so that these are used instead of Go's default HTTP server settings
	srv := http.Server{
		Addr:      *addr,
		Handler:   mux,
		ErrorLog:  errLog,
		TLSConfig: &tlsConfig,
		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting Server on %s", *addr)

	//call the listen and serve TLS method on our new htpp server
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	if err != nil {
		errLog.Fatal(err)
	}
}
