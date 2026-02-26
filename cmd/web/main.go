package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	defaultAddr := os.Getenv("SNIPPETBOX_ADDR")
	if defaultAddr == "" {
		defaultAddr = "127.0.0.1:4000"
	}

	addr := flag.String("addr", defaultAddr, "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "\033[32mINFO\t\033[0m", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "\033[31mERROR\t\033[0m", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:   app.routes(),
	}


	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	if err != nil {
		errorLog.Fatalf("Server failed to start: %v", err)
	}
}
