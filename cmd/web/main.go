package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ptakpatryk/list-me/internals/models"
)

type application struct {
	logger        *slog.Logger
	lists         *models.ListModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
  sessionManager *scs.SessionManager
}

func main() {
	connectionString := flag.String("db", "", "postgresql connection string")
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	db, err := openDB(*connectionString)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

  sessionManager := scs.New()
  sessionManager.Store = postgresstore.New(db)
  sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger:        logger,
		lists:         &models.ListModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
    sessionManager: sessionManager,
	}

	logger.Info("starting server", slog.String("addr", *addr))

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
