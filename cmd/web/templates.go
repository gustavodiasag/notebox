package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/gustavodiasag/notebox/internal/models"
)

// Acts as the holding structure for any dynamic data passed to HTML templates.
type templateData struct {
	CurrentYear     int
	Note            *models.Note
	Notes           []*models.Note
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func fmtDate(t time.Time) string {
	return t.Format("02 Jan, 2006")
}

var functions = template.FuncMap{
	"fmtDate": fmtDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as a cache.
	cache := map[string]*template.Template{}

	// Provides a slice of all the filepaths for the application templates.
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
