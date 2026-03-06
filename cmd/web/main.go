package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	"github.com/iancenry/snippetbox/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder *form.Decoder
	sessionManager *scs.SessionManager
}



func main() {
	godotenv.Load()
	// parses runtime config settings for the app
	defaultAddr := os.Getenv("SNIPPETBOX_ADDR")
	if defaultAddr == "" {
		defaultAddr = "127.0.0.1:4000"
	}

	addr := flag.String("addr", defaultAddr, "HTTP network address")
	dsn := flag.String("dsn", os.Getenv("DATABASE_URL"), "PostgreSQL connection string")
	flag.Parse()

	infoLog := log.New(os.Stdout, "\033[32mINFO\t\033[0m", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "\033[31mERROR\t\033[0m", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)

	// establishes a connection pool to the database
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatalf("Unable to connect to database: %v", err)
	}
	infoLog.Println("Database connection pool established")
	defer db.Close()

	// initializes a new session manager and configures it to use the PostgreSQL database as the session store
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true


	
	// establishes dependencies for handlers - loggers and database models
	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		snippets: &models.SnippetModel{DB: db},
		formDecoder: form.NewDecoder(),
		sessionManager: sessionManager,
	}

	// create a template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatalf("Unable to create template cache: %v", err)
	}
	app.templateCache = templateCache

	// runs the http server and listens for requests
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:   app.routes(),
	}


	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatalf("Server failed to start: %v", err)
	}
}


func openDB(dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return pool, nil
}