# System Architecture - Personal Finance CLI Manager

## Overview

The application is a CLI-based personal finance manager. The system is designed with modularity and testability in mind, following Go best practices.

## Layers

### 1. CLI Layer

* **Location:** `internal/cli/`
* **Entry Point:** `cmd/main.go`
* **Responsibility:** Command parsing using Cobra, input validation, and terminal output formatting.

### 2. Data Layer

* **Location:** `internal/db/`
* **Responsibility:** SQLite database connection, migrations, persistent storage.
* **Migrations:** SQL scripts in `internal/db/migrations/`
* **Connection:** `db.go` provides `NewDB()` for opening and configuring SQLite.

### 3. Models Layer

* **Location:** `internal/models/`
* **Responsibility:** Domain models + CRUD methods.
* **Models Implemented:**

    * `Transaction` (id, date, description, amount, category)
    * `Budget` (id, category, amount, period)
    * `CategoryRule` (id, regex pattern, category)

### 4. Testing Layer

* **Location:** `tests/`
* **Responsibility:** Unit tests for models and DB
* **Test DB:** in-memory SQLite (`:memory:`)
* **Migration Loader:** Automatically applies all migration scripts before tests

### 5. Future Layers

* Importers (CSV/OFX parsing)
* Categorization Engine (regex-based auto-categorization)
* Reporting & Budget alerts
* TUI (BubbleTea)

## Data Flow

1. CLI parses user input → calls model/service functions
2. Model methods interact with SQLite via `internal/db` connection
3. Transactions, Budgets, and Rules stored in database tables
4. Reports and category assignments use model methods to fetch and process data

## Folder Structure Summary

```
personal-finance-cli/
├── cmd/              # CLI entry point
├── internal/
│   ├── db/           # Database connection + migrations
│   └── models/       # Models and CRUD
└── tests/            # Unit tests
```

## Design Principles

* **Modularity:** each layer separated for clarity and testability
* **Testability:** in-memory DB, migration loader, isolated unit tests
* **Extensibility:** easy to add new CLI commands, importers, or TUI components
