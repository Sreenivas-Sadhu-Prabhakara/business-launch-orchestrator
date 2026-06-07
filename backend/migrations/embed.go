// Package migrations embeds the SQL migration files so they ship inside the
// compiled binary and run automatically on startup — no external migration
// tool required.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
