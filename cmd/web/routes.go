package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	router.Handler(http.MethodGet,"/static/*filepath", http.StripPrefix("/static/", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView) 
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)
	router.HandlerFunc(http.MethodGet, "/snippets", app.snippetLatest )

	// standard middleware chain for all requests
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}


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

