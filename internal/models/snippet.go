package models

import (
	"time"
)

type (
	Snippet struct {
		Created time.Time
		Expires time.Time
		Title   string
		Content string
		ID      int
	}
)
