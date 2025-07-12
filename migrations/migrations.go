package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var migrationFS embed.FS

// Applies all migrations to the database.
func Up(db *sql.DB) error {
	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(migrationFS)
	err := goose.SetDialect("pgx")
	if err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}
	err = goose.Up(db, ".")
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}
