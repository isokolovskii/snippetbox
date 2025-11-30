package models

import (
	"errors"
)

// ErrNoRecord - error returned if snippet not found in database.
var ErrNoRecord = errors.New("models: no matching record found")
