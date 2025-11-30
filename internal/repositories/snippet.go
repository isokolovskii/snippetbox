package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"snippetbox.isokol.dev/internal/models"
)

type (
	// SnippetRepository - snippets repository for snippet model.
	SnippetRepository struct {
		// Database connection.
		db *sql.DB
	}
)

const (
	// SQL query for snippet insertion.
	snippetInsertQuery = `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	// SQL query part for select fields on snippets.
	snippetSelectQueryPart = "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP()"
	// SQL query for snippet get.
	snippetGetQueryPart = " AND id = ?"
	// SQL query for latest 10 snippets.
	snippetLatestQueryPart = "ORDER BY id DESC LIMIT 10"
)

// Insert - insert snippet into database.
func (m *SnippetRepository) Insert(ctx context.Context, title, content string, expires int) (int, error) {
	result, err := m.db.ExecContext(ctx, snippetInsertQuery, title, content, expires)
	if err != nil {
		return 0, fmt.Errorf("error inserting new snippet into database: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting ID of last inserted element: %w", err)
	}

	return int(id), nil
}

// Get snippet by ID from database.
func (m *SnippetRepository) Get(ctx context.Context, id int) (models.Snippet, error) {
	row := m.db.QueryRowContext(ctx, snippetSelectQueryPart+snippetGetQueryPart, id)

	var snippet models.Snippet

	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Snippet{}, models.ErrNoRecord
		}

		return models.Snippet{}, fmt.Errorf("error during search of snippet by ID: %w", err)
	}

	return snippet, nil
}

// Latest - get latest snippets from database.
func (m *SnippetRepository) Latest(ctx context.Context) ([]models.Snippet, error) {
	rows, err := m.db.QueryContext(ctx, snippetSelectQueryPart+snippetLatestQueryPart)
	if err != nil {
		return nil, fmt.Errorf("error querying latest snippets from database: %w", err)
	}
	defer rows.Close()

	var snippets []models.Snippet

	for rows.Next() {
		var snippet models.Snippet

		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, fmt.Errorf("error creating models from queried database rows: %w", err)
		}
		snippets = append(snippets, snippet)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error selecting latest snippets from database: %w", err)
	}

	return snippets, nil
}
