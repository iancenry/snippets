package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/iancenry/snippetbox/internal/models"
	"github.com/iancenry/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type SnippetCreateForm struct {
	Title string  `form:"title"`
	Content string `form:"content"`
	Expires int `form:"expires"`
	validator.Validator `form:"-"`
}


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


	id, err := uuid.Parse(params.ByName("id"))
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

	flash := app.sessionManager.PopString(r.Context(), "flash")
	
	data := app.newTemplateData(r)
	data.Snippet = snippet
	data.Flash = flash

	app.render(w, http.StatusOK, "view.tmpl.html", data)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request){

	data := app.newTemplateData(r)
	data.Form = SnippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl.html", data)
}


func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request){
	
	var form SnippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.Check(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.Check(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.Check(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.Check(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%s", id.String()), http.StatusSeeOther)
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