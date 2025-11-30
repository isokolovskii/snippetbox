package migrations

import (
	"embed"
)

// Files - embedded database migrations.
//
//go:embed "*.sql"
var Files embed.FS
