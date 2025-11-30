package main

import (
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.isokol.dev/internal/models"
)

type (
	// Template data.
	templateData struct {
		// Message for flash messaging.
		Flash string
		// Snippet entity.
		Snippet *models.Snippet
		// Form for forms refill after error.
		Form any
		// Snippets array for showing multiple snippets as list.
		Snippets []models.Snippet
		// Current year.
		CurrentYear int
	}
)

// Create template cache.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, fmt.Errorf("failed to read all pages templates: %w", err)
	}

	for _, page := range pages {
		name := filepath.Base(page)

		tmpl, err := parsePageTemplate(name, page)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template: %w", err)
		}

		cache[name] = tmpl
	}

	return cache, nil
}

// Parsing of template html.
func parsePageTemplate(name, page string) (*template.Template, error) {
	functions := template.FuncMap{
		"humanDate": humanDate,
	}
	tmpl, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse base template: %w", err)
	}

	tmpl, err = tmpl.ParseGlob("./ui/html/partials/*.tmpl.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates for partials: %w", err)
	}

	tmpl, err = tmpl.ParseFiles(page)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	return tmpl, nil
}

// Date formatting for template rendering.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
