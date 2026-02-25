package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// neuteredFileSystem is a custom http.FileSystem that returns os.ErrNotExist
// for any directories, preventing directory listing. - alexedwards.net/blog/disable-http-fileserver-directory-listings
type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, os.ErrNotExist
		}
	}

	return f, nil
}

func main() {
	defaultAddr := os.Getenv("SNIPPETBOX_ADDR")
	if defaultAddr == "" {
		defaultAddr = "127.0.0.1:4000"
	}

	addr := flag.String("addr", defaultAddr, "HTTP network address")
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))


	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view/", snippetView) 
	mux.HandleFunc("/snippet/create", snippetCreate)



	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
