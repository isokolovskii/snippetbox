package models

import (
	"errors"
)

var (
	// ErrNoRecord - error returned if snippet not found in database.
	ErrNoRecord = errors.New("models: no matching record found")
	// ErrInvalidCredentials - error returned if user provided wrong credentials
	// during login.
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// ErrDuplicateEmail - error returned if user with specified email
	// already exists in database and insert operation fails due to duplicate.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
