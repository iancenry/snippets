package main

import (
	"html/template"
	"net/http"
	"strconv"
)


func (app *application) home(w  http.ResponseWriter, r *http.Request){
	if(r.URL.Path != "/"){
		http.NotFound(w, r)
		return
	}

	if(r.Method != http.MethodGet){
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Printf("Error parsing template files: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.errorLog.Printf("Error executing template: %v", err)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// w.Write([]byte("Hello from snippetbox"))
}

func (app *application) snippetView(w  http.ResponseWriter, r *http.Request){
	if(r.Method != http.MethodGet){
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	stringId := r.URL.Query().Get("id")

	if id, err := strconv.Atoi(stringId); err != nil || id < 1 {
		app.errorLog.Printf("Invalid snippet ID: %v", stringId)
		
		http.NotFound(w, r)
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
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}


	w.Write([]byte("Create a new snippet"))
}