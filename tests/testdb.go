package tests

import (
	"database/sql"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	_ "modernc.org/sqlite"
)

func NewTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	applyMigrations(t, db)

	return db
}

func applyMigrations(t *testing.T, db *sql.DB) {
	t.Helper()

	// Get absolute path to this file
	_, filename, _, _ := runtime.Caller(0)

	// migrations folder relative to testdb.go
	migrationsDir := filepath.Join(filepath.Dir(filename), "..", "internal", "db", "migrations")

	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("failed to read migrations directory: %v", err)
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".sql" {
			continue
		}

		content, err := ioutil.ReadFile(filepath.Join(migrationsDir, f.Name()))
		if err != nil {
			t.Fatalf("failed to read migration file %s: %v", f.Name(), err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			t.Fatalf("migration %s failed: %v", f.Name(), err)
		}
	}
}
