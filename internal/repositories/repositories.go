package repositories

import (
	"database/sql"
)

type (
	Repositories struct {
		Snippet *SnippetRepository
	}
)

func CreateRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Snippet: &SnippetRepository{
			db: db,
		},
	}
}
