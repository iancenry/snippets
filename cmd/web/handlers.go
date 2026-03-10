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


type UserSignupForm struct {
	Name string `form:"name"`
	Email string `form:"email"`
	Password string `form:"password"`
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
	
	data := app.newTemplateData(r)
	data.Snippet = snippet

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
	form.Check(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

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

func (app *application) userSignup(w http.ResponseWriter, r *http.Request){
	data := app.newTemplateData(r)
	data.Form = UserSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request){
	var form UserSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.Check(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.Check(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.Check(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.Check(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.Check(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	_, err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddError("email", "Address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request){
	data := app.newTemplateData(r)
	data.Form = UserSignupForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request){
	var form UserSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.Check(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.Check(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.Check(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request){
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping (w http.ResponseWriter, r *http.Request){
	w.Write([]byte("OK"))
}