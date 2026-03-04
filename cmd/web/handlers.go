package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/iancenry/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
)


func (app *application) home(w  http.ResponseWriter, r *http.Request){
	// panic("Oh no, a problem occurred!")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)

	// w.Write([]byte("Hello from snippetbox"))
}

func (app *application) snippetView(w  http.ResponseWriter, r *http.Request){
	params := httprouter.ParamsFromContext(r.Context())

	stringId := params.ByName("id")

	id, err := uuid.Parse(stringId)
	if err != nil {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request){

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"message": "Created a new snippet with id: " + id.String()})
}


func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request){

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"message": "Created a new snippet with id: " + id.String()})
}


func (app *application) snippetLatest(w http.ResponseWriter, r *http.Request){
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(snippets)
}