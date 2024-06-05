package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/ptakpatryk/list-me/internals/models"
	"github.com/ptakpatryk/list-me/ui"
)

type templateData struct {
	CurrentYear     int
	List            models.List
	Lists           []models.List
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}
  

	for _, page := range pages {
		name := filepath.Base(page)

    patterns := []string{
      "html/base.tmpl.html",
      "html/partials/*.tmpl.html",
      page,
    }

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
