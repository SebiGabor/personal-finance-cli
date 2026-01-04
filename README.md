# Personal Finance CLI Manager

A command-line tool for tracking personal income and expenses. This project is being developed step-by-step as a complete Go application.

## Current Status

* Go module initialized
* SQLite connection established
* Database migrations implemented
* Core models: Transaction, Budget, CategoryRule
* CRUD methods implemented for all models
* Unit tests using in-memory SQLite
* CLI entry point placeholder

## Project Structure

```
personal-finance-cli/
│   go.mod
│   go.sum
│   README.md
│
├── cmd/
│   └── main.go
│
├── internal/
│   ├── cli/
│   │   ├── add.go
│   │   ├── list.go
│   │   └── root.go
│   ├── db/
│   │   ├── db.go
│   │   └── migrations/
│   │       ├── 001_init.sql
│   │       └── 002_category_rules.sql
│   │
│   └── models/
│       ├── transaction.go
│       ├── budget.go
│       └── category_rule.go
│
└── tests/
    ├── testdb.go
    ├── transaction_test.go
    ├── budget_test.go
    ├── cli_test.go
    └── category_rule_test.go
```

## Database Layer

* Opens SQLite DB (`finance.db` in production, `:memory:` for tests)
* Foreign keys enabled
* Migrations automatically applied in tests

## Models Implemented

* **Transaction**: CRUD, date, description, amount, category
* **Budget**: CRUD, category, amount, period
* **CategoryRule**: CRUD, regex-based categorization rules

## Testing

* Uses in-memory SQLite
* Migration loader ensures DB schema is correct
* Tests cover creation, reading, updating, deletion

## Running the Project

```
go run ./cmd
```

## Running Tests

```
go test ./...
```

## Next Steps

* Auto-run migrations on startup
* Cobra CLI structure (`add`, `import`, `budget`, `search`, `report`)
* CSV/OFX importers
* Regex-based auto-categorizer
* Budget alerts
* Reports and terminal charts
* Interactive TUI (BubbleTea)
