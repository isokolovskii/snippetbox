package mocks

import (
	"context"
	"time"

	"snippetbox.isokol.dev/internal/models"
)

type (
	SnippetRepository struct{}
)

var mockSnippet = models.Snippet{
	ID:      1,
	Title:   "Mock snippet",
	Content: "I am mock snippet",
	Created: time.Now(),
	Expires: time.Now(),
}

func (*SnippetRepository) Insert(_ context.Context, _, _ string, _ int) (int, error) {
	return 2, nil
}

func (*SnippetRepository) Get(_ context.Context, id int) (models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return models.Snippet{}, models.ErrNoRecord
	}
}

func (*SnippetRepository) Latest(_ context.Context) ([]models.Snippet, error) {
	return []models.Snippet{mockSnippet}, nil
}
