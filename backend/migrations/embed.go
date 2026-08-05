// Package migrations embeds the SQL migration files so they can be applied by
// the golang-migrate iofs source driver without shipping the .sql files.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
