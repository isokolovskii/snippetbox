package main

import (
	"fmt"
	"path/filepath"
	"text/template"

	"snippetbox.isokol.dev/internal/models"
)

type (
	templateData struct {
		Snippet     models.Snippet
		Snippets    []models.Snippet
		CurrentYear int
	}
)

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

func parsePageTemplate(name, page string) (*template.Template, error) {
	tmpl, err := template.ParseFiles("./ui/html/base.tmpl.html")
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
