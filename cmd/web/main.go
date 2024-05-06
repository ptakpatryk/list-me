package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

  //Static server
  fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

  //Handlers
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /list/view/{id}", listView)
	mux.HandleFunc("GET /list/create", listCreate)
	mux.HandleFunc("POST /list/create", listCreatePost)


	log.Print("starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
