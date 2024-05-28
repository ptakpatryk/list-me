package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	//Static server
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	//Handlers
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /list/view/{id}", dynamic.ThenFunc(app.listView))
	mux.Handle("GET /list/create", dynamic.ThenFunc(app.listCreate))
	mux.Handle("POST /list/create", dynamic.ThenFunc(app.listCreatePost))

  //Auth
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", dynamic.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
