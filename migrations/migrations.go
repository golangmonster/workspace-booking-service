package migrations

import "embed"

// Files contains the same SQL migrations used by the application.
//
//go:embed *.sql
var Files embed.FS
