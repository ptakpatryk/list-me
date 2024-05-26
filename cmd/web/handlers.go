package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ptakpatryk/list-me/internals/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	lists, err := app.lists.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.Lists = lists

	app.render(w, r, 200, "home.tmpl.html", data)
}

func (app *application) listView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	list, err := app.lists.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.List = list

	app.render(w, r, 200, "view.tmpl.html", data)
}

func (app *application) listCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

type snippetCreateForm struct {
	Title       string
	Description string
	Expires     int
	FieldErrors map[string]string
}

func (app *application) listCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Description: r.PostForm.Get("description"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form.Description) == "" {
		form.FieldErrors["description"] = "This field cannot be blank"
	}

	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.lists.Insert(form.Title, form.Description, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/list/view/%d", id), http.StatusSeeOther)
}
