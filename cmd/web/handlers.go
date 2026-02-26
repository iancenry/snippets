package main

import (
	"html/template"
	"net/http"
	"strconv"
)


func (app *application) home(w  http.ResponseWriter, r *http.Request){
	if(r.URL.Path != "/"){
		app.notFound(w)
		return
	}

	if(r.Method != http.MethodGet){
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// w.Write([]byte("Hello from snippetbox"))
}

func (app *application) snippetView(w  http.ResponseWriter, r *http.Request){
	if(r.Method != http.MethodGet){
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	stringId := r.URL.Query().Get("id")

	if id, err := strconv.Atoi(stringId); err != nil || id < 1 {
		app.notFound(w)

		return
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`{"message": "Display a specific snippet with id: ` + strconv.Itoa(id) + `"}`))
		return	
	}
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request){
	if(r.Method != http.MethodPost){
		w.Header().Set("Allow", http.MethodPost)

		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}


	w.Write([]byte("Create a new snippet"))
}