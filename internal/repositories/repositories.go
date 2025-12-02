package repositories

import (
	"database/sql"
)

type (
	// Repositories - database repositories.
	Repositories struct {
		// Snippets repository.
		Snippet SnippetRepositoryInterface
		// Users repository.
		User UserRepositoryInterface
	}
)

// CreateRepositories - create repositories.
func CreateRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Snippet: &SnippetRepository{
			db: db,
		},
		User: &UserRepository{
			db: db,
		},
	}
}
