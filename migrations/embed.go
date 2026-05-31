// Package migrations встраивает SQL-файлы миграций в бинарник.
package migrations

import "embed"

// FS содержит все .sql файлы миграций.
//
//go:embed *.sql
var FS embed.FS
