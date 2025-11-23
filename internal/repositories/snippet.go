package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"snippetbox.isokol.dev/internal/models"
)

type (
	SnippetRepository struct {
		db *sql.DB
	}
)

func (m *SnippetRepository) Insert(ctx context.Context, title, content string, expires int) (int, error) {
	query := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.db.ExecContext(ctx, query, title, content, expires)
	if err != nil {
		return 0, fmt.Errorf("error inserting new snippet into database: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting ID of last inserted element: %w", err)
	}

	return int(id), nil
}

func (m *SnippetRepository) Get(ctx context.Context, id int) (models.Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.db.QueryRowContext(ctx, query, id)

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

func (m *SnippetRepository) Latest(ctx context.Context) ([]models.Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.db.QueryContext(ctx, query)
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
