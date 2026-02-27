package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
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

	// establishes dependencies for handlers - loggers
	infoLog := log.New(os.Stdout, "\033[32mINFO\t\033[0m", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "\033[31mERROR\t\033[0m", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	// establishes a connection pool to the database
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatalf("Unable to connect to database: %v", err)
	}
	infoLog.Println("Database connection pool established")
	defer db.Close()



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