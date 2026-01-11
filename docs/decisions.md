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

## 9. CSV Import Strategy

* **Decision:** Use Go's standard `encoding/csv` library.
* **Format:** Enforce a simple default CSV structure (`Date,Description,Amount,Category`) for the initial version.
* **Reason:** Provides immediate value with minimal dependencies. Allows users to bulk-import data by converting their bank statements to this simple format before we implement complex bank-specific parsers.

## 10. Reporting Strategy

* **Decision:** Implement monthly aggregation using SQL `GROUP BY` and display simple ASCII bar charts in the terminal.
* **Reason:** Provides immediate visual feedback and spending insights without requiring external UI libraries or complex dependencies. Fits the "CLI-first" philosophy.

## 11. Auto-Categorization Logic

* **Decision:** Implement a rule-based categorization engine using Go's `regexp` library.
* **Reason:** Enables flexible, case-insensitive pattern matching (e.g., `(?i)uber`) to automatically assign categories to transactions during import, significantly reducing manual data entry.

## 12. Budget Management

* **Decision:** Create a dedicated `budget` command with `add`, `list`, and `remove` subcommands.
* **Reason:** Separates configuration (setting limits) from reporting (viewing progress). Keeps the CLI organized and modular. Defaults to "monthly" periods for simplicity in the initial version.

## 13. Budget Alerts and Monitoring

* **Decision:** Implement active monitoring during transaction creation (`add` command) and passive monitoring via visual indicators in the `budget list` command.
* **Reason:**
    * **Visuals:** ASCII progress bars provide immediate "at a glance" status of finances without needing complex GUI charts.
    * **Alerts:** Triggering alerts immediately after a manual entry gives the user instant feedback on their spending habits.
    * **Non-blocking:** We chose *not* to block the transaction if it exceeds the budget, as the CLI is a tracking tool, not a bank authorization system. The data must reflect reality.

## 14. Interactive Terminal UI (TUI)

* **Decision:** Use the `github.com/rivo/tview` library to implement an interactive browsing mode.
* **Reason:**
    * **User Experience:** CLI flags (`list`, `search`) are great for quick queries, but browsing a large history of transactions is better served by a scrollable visual table.
    * **Library Choice:** `tview` was chosen over `termui` or raw `tcell` because it provides high-level components (Tables, Forms) that simplify rendering and event handling, reducing development time.

## 15. OFX File Support

* **Decision:** Implement a custom XML parser for Open Financial Exchange (OFX) files using Go's `encoding/xml`.
* **Reason:**
    * **Standardization:** OFX is a widely used standard by banks for exporting transaction data, offering a more reliable structure than CSV (which varies by bank).
    * **Simplicity:** Instead of a heavy third-party library, we defined minimal Go structs matching the specific XML tags needed (`<STMTTRN>`, `<TRNAMT>`, etc.) to keep dependencies low.