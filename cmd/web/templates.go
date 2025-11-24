package main

import (
	"snippetbox.isokol.dev/internal/models"
)

type (
	templateData struct {
		Snippet  models.Snippet
		Snippets []models.Snippet
	}
)
