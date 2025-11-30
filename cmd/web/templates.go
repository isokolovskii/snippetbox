package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.isokol.dev/internal/models"
	"snippetbox.isokol.dev/ui"
)

type (
	// Template data.
	templateData struct {
		// Snippet entity.
		Snippet *models.Snippet
		// Snippets array for showing multiple snippets as list.
		Snippets []models.Snippet
		// Form for forms refill after error.
		Form any
		// Message for flash messaging.
		Flash string
		// Current year.
		CurrentYear int
		// User authentication status.
		IsAuthenticated bool
		// CRSF token.
		CSRFToken string
	}
)

// Create template cache.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
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

	patterns := []string{
		"html/base.tmpl.html",
		"html/partials/*.tmpl.html",
		page,
	}

	tmpl, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return tmpl, nil
}

// Date formatting for template rendering.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
