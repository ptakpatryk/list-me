package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
  addr := flag.String("addr", ":4000", "HTTP network address")
  flag.Parse()

  logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    AddSource: true,
    Level: slog.LevelDebug,
  }))

	mux := http.NewServeMux()

  //Static server
  fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

  //Handlers
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /list/view/{id}", listView)
	mux.HandleFunc("GET /list/create", listCreate)
	mux.HandleFunc("POST /list/create", listCreatePost)


  logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, mux)
  logger.Error(err.Error())
	os.Exit(1)
}
