# Design Decisions - Personal Finance CLI Manager

## 1. Go Module Structure

* **Decision:** Use `cmd/`, `internal/`, `tests/` structure
* **Reason:** Standard Go project layout, keeps internal packages private, separates executable code from reusable modules
* **Reference:** [https://blog.golang.org/organizing-go-code](https://blog.golang.org/organizing-go-code)

## 2. Database

* **Decision:** SQLite for local storage
* **Reason:** Lightweight, zero-configuration, cross-platform, perfect for CLI apps
* **Decision:** Use `modernc.org/sqlite` instead of `github.com/mattn/go-sqlite3`
* **Reason:** Pure Go, easier on Windows (no cgo dependency)

## 3. Models Layer

* **Decision:** Implement Transaction, Budget, CategoryRule structs with CRUD methods
* **Reason:** Encapsulates DB operations, provides clear domain representation
* **Decision:** Pass `*sql.DB` to all methods instead of global DB instance
* **Reason:** Better testability and separation of concerns

## 4. Migrations

* **Decision:** SQL files in `internal/db/migrations/`
* **Reason:** Simple, readable, portable; allows in-memory DB setup for tests
* **Decision:** Load migrations automatically in test DB using helper
* **Reason:** Ensure DB schema is always consistent before tests

## 5. Testing

* **Decision:** In-memory SQLite for unit tests
* **Reason:** Fast, no files created, repeatable tests
* **Decision:** Test all CRUD operations for all models
* **Reason:** Ensures DB schema + logic correctness early in development

## 6. CLI & Future Decisions

* **Decision:** Plan to use Cobra for subcommands (`add`, `import`, `budget`, `report`, `search`)
* **Reason:** Standard library for building CLI apps in Go, supports subcommands, flags, and interactive prompts
* **Decision:** Use regex-based categorization rules
* **Reason:** Simple yet flexible auto-categorization engine
* **Decision:** Plan TUI using BubbleTea
* **Reason:** Rich terminal UI for browsing and reporting transactions

## 7. CLI Package Separation

* **Decision:** Move Cobra command logic from `cmd/` to `internal/cli/`.
* **Reason:** Keeps the main entry point clean and allows the CLI logic to be tested by injecting a mock or test database via an exported `SetDatabase` function.

## 8. Transaction Search and Filter

* **Decision:** Implement search using the SQL `LIKE` operator within the models layer.
* **Reason:** Keeping the filtering logic in SQL is more efficient than fetching all records and filtering in Go. It fulfills the user story for searching and filtering with minimal complexity.