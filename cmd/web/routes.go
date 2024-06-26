package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/ptakpatryk/list-me/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	//Static server
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

  mux.HandleFunc("GET /ping", ping)

	//Handlers
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /list/view/{id}", dynamic.ThenFunc(app.listView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected
	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET /list/create", protected.ThenFunc(app.listCreate))
	mux.Handle("POST /list/create", protected.ThenFunc(app.listCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
