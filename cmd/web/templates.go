package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/gustavodiasag/notebox/internal/models"
	"github.com/gustavodiasag/notebox/ui"
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
