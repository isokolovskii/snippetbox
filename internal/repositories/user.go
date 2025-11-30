package repositories

import (
	"database/sql"
)

type (
	// UserRepository - database repository for user model.
	UserRepository struct {
		// Database connection.
		db *sql.DB
	}
)

// Insert - insert new user to database.
func (*UserRepository) Insert(_, _, _ string) error {
	return nil
}

// Authenticate - verify user credentials using provided
// email and password.
func (*UserRepository) Authenticate(_, _ string) (int, error) {
	return 0, nil
}

// Exists - check if user exists in database.
func (*UserRepository) Exists(_ int) (bool, error) {
	return false, nil
}
