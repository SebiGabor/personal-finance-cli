# Personal Finance CLI Manager

A robust, local-first command-line tool for tracking personal income and expenses. Built with Go, it allows you to import transactions from bank statements (CSV/OFX), categorize them automatically using regex rules, set monthly budgets, and visualize your spending habits—all directly from your terminal.

![Project Status](https://img.shields.io/badge/status-complete-success)
![Go Version](https://img.shields.io/badge/go-1.25%2B-blue)

## Features

* **Multi-Format Import:** Seamlessly import transactions from **CSV** and **OFX** (Bank Export) files.
* **Auto-Categorization:** Define Regex-based rules to automatically assign categories to new transactions.
* **Duplicate Detection:** Smart import logic prevents duplicate entries, even if you re-import the same file.
* **Budgeting & Alerts:** Set monthly limits per category. The CLI warns you immediately if you overspend.
* **Visual Reports:** Generate ASCII bar charts to visualize monthly spending breakdowns.
* **Interactive TUI:** Browse, scroll, and view your transaction history in a rich Terminal UI.
* **Search & Filter:** Instantly find transactions by keyword or category.
* **Local Storage:** Uses a zero-dependency SQLite database (no server required).

---

## Installation

### Prerequisites
* [Go 1.25+](https://go.dev/dl/) installed.

### Build from Source
1. **Clone the repository:**
   ```bash
   git clone https://github.com/SebiGabor/personal-finance-cli.git
   cd personal-finance-cli
   ```

2. **Build the executable:**
   ```bash
   go build -o finance.exe ./cmd
   ```
   *This creates a binary file named `finance` (or `finance.exe` on Windows) in the current folder.*

---

## Usage Guide

To run the tool, use the `./finance` command followed by a subcommand.

### 1. Import Transactions
Import data from external files. The system auto-normalizes categories (e.g., "food" -> "Food") and skips duplicates.

```bash
# Import a CSV file
./finance import test.csv

# Import an OFX file (Standard Bank Export)
./finance import test.ofx
```

### 2. Manual Entry
Add cash expenses or income manually.

```bash
# Add an expense (negative amount)
./finance add --amount -15.50 --desc "Lunch at McD" --category "Food"

# Add income (positive amount)
./finance add --amount 2500 --desc "Freelance Project" --category "Income" --date 2024-02-01
```

### 3. Budget Management
Set spending limits to keep your finances in check.

```bash
# Set a $500 monthly limit for Groceries
./finance budget add --category "Groceries" --amount 500

# View all budgets and current progress
./finance budget list

# Remove a budget (find ID via 'list')
./finance budget remove [ID]
```

### 4. Reporting
Visualize your financial health with ASCII charts.

```bash
# Show spending breakdown for the current month
./finance report

# Show report for a specific month/year
./finance report --month 12 --year 2023
```
*Output Example:*
```text
Monthly Report for January 2024
--------------------------------------------------
Category Breakdown:
Groceries       │ ████████████ 450.00
Entertainment   │ ████ 120.00
Transport       │ ██ 60.00
```

### 5. Interactive Mode (TUI)
Launch the visual interface to browse your full transaction history.
* **Navigation:** Use `Arrow Keys` to scroll up/down.
* **Quit:** Press `q` or `Esc` to exit.

```bash
./finance tui
```

### 6. Search
Find specific transactions by description or category.

```bash
./finance search "Uber"
./finance search "Salary"
```

### 7. Automation Rules
Manage regex rules to auto-categorize future imports.

```bash
# Map any transaction containing "Netflix" to "Entertainment"
./finance rules add --pattern "Netflix" --category "Entertainment"

# List all active rules
./finance rules list

# Remove a rule
./finance rules remove [ID]
```

### 8. Manage Transactions
View or delete individual transactions.

```bash
# List all transactions (newest first)
./finance list

# Delete a specific transaction (find ID via 'list' or 'search')
./finance delete [ID]
```

---

## Project Structure

The project follows the standard [Golang Project Layout](https://github.com/golang-standards/project-layout):

```text
personal-finance-cli/
├── cmd/                # Entry point (main.go)
├── internal/
│   ├── cli/            # Command logic (Cobra handlers)
│   ├── db/             # SQLite connection & Migrations
│   ├── models/         # Data structures & Business logic
│   └── tui/            # Terminal UI implementation (tview)
├── tests/              # Integration tests
├── docs/               # Architecture & Decision logs
├── go.mod              # Dependency definitions
└── README.md           # This file
```

## Testing

The project uses a comprehensive suite of integration tests with an in-memory database to ensure reliability.

```bash
# Run all tests
go test ./tests/...
```