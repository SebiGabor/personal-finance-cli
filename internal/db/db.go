package db

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Connect opens (or creates) the SQLite database file and runs migrations.
func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "finance.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// runMigrations executes all .sql files in the migrations folder.
func runMigrations(db *sql.DB) error {
	entries, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		migrationName := entry.Name()
		content, err := migrationFiles.ReadFile("migrations/" + migrationName)
		if err != nil {
			return err
		}

		log.Printf("Running migration: %s", migrationName)

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("migration %s failed: %w", migrationName, err)
		}
	}

	return nil
}
