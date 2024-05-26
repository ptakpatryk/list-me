package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) listCreatePost(w http.ResponseWriter, r *http.Request) {
  err := r.ParseForm()
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }


	title := r.PostForm.Get("title")
	description := r.PostForm.Get("description")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

	id, err := app.lists.Insert(title, description, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/list/view/%d", id), http.StatusSeeOther)
}
