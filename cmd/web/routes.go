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

	//Handlers
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /list/view/{id}", app.listView)
	mux.HandleFunc("GET /list/create", app.listCreate)
	mux.HandleFunc("POST /list/create", app.listCreatePost)

  standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
