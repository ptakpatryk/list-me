package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ptakpatryk/list-me/internals/models"
	"github.com/ptakpatryk/list-me/internals/validator"
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

type listCreateForm struct {
	Title               string `form:"title"`
	Description         string `form:"description"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) listCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = listCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) listCreatePost(w http.ResponseWriter, r *http.Request) {
	var form listCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Description), "description", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
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

	app.sessionManager.Put(r.Context(), "flash", "Snippet succesfully created!")

	http.Redirect(w, r, fmt.Sprintf("/list/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Signup form")
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request)  {
  fmt.Fprintln(w, "Create a new user")
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Login form")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Authenticate and login the user")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Logout the user...")
}
