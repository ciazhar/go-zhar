package migrations

import "embed"

//go:embed *.json
var MigrationsFS embed.FS
