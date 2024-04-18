package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Bude8/lets-go/internal/models"
)

// html/template only allows one item of dynamic data when rendering a template
// This allows multiple pieces of dynamic data
type templateData struct {
	Snippet     models.Snippet
	Snippets    []models.Snippet
	CurrentYear int
	Form        any
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
	}
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialise a new map to act as the cache
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Register template.FuncMap with template set before ParseFiles.
		// Use template.New() to create an empty ts, then register, then parse.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
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
