package main

import (
	"html/template"
	"path/filepath"

	"github.com/gustavodiasag/notebox/internal/models"
)

// Acts as the holding structure for any dynamic data passe to HTML templates.
type templateData struct {
	Note  *models.Note
	Notes []*models.Note
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

		files := []string{
			"./ui/html/base.tmpl.html",
			"./ui/html/partials/nav.tmpl.html",
			page,
		}
		// Parse the file into a template set.
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
