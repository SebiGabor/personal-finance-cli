package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Budget struct {
	ID       int64
	Category string
	Amount   float64
	Period   string // "monthly", "weekly", "yearly"
}

// CreateBudget creates a new budget or updates an existing one for the category
func CreateBudget(db *sql.DB, b *Budget) error {
	// 1. Check if a budget for this category already exists
	var existingID int64 // FIXED: Changed from int to int64
	err := db.QueryRow("SELECT id FROM budgets WHERE category = ?", b.Category).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Case A: No budget exists, create a new one (INSERT)
		query := `INSERT INTO budgets (category, amount, period) VALUES (?, ?, ?)`
		result, err := db.Exec(query, b.Category, b.Amount, b.Period)
		if err != nil {
			return fmt.Errorf("failed to insert budget: %w", err)
		}
		id, _ := result.LastInsertId() // LastInsertId returns int64
		b.ID = id                      // FIXED: No longer casting to int(id)
	} else if err != nil {
		// Database error
		return fmt.Errorf("failed to check existing budget: %w", err)
	} else {
		// Case B: Budget exists, update it (UPDATE)
		query := `UPDATE budgets SET amount = ?, period = ? WHERE id = ?`
		_, err := db.Exec(query, b.Amount, b.Period, existingID)
		if err != nil {
			return fmt.Errorf("failed to update budget: %w", err)
		}
		b.ID = existingID // FIXED: existingID is now int64, so this assignment works
	}

	return nil
}

func GetBudget(db *sql.DB, id int64) (*Budget, error) {
	query := `
        SELECT id, category, amount, period
        FROM budgets WHERE id = ?;
    `
	row := db.QueryRow(query, id)

	var b Budget
	if err := row.Scan(&b.ID, &b.Category, &b.Amount, &b.Period); err != nil {
		return nil, err
	}

	return &b, nil
}

func ListBudgets(db *sql.DB) ([]Budget, error) {
	rows, err := db.Query("SELECT id, category, amount, period FROM budgets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []Budget
	for rows.Next() {
		var b Budget
		if err := rows.Scan(&b.ID, &b.Category, &b.Amount, &b.Period); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}
	return budgets, nil
}

func UpdateBudget(db *sql.DB, b *Budget) error {
	_, err := db.Exec(`
        UPDATE budgets
        SET category = ?, amount = ?, period = ?
        WHERE id = ?;
    `, b.Category, b.Amount, b.Period, b.ID)
	return err
}

func DeleteBudget(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM budgets WHERE id = ?", id)
	return err
}

func GetSpendingTotal(db *sql.DB, category string, month time.Month, year int) (float64, error) {
	dateFilter := fmt.Sprintf("%04d-%02d%%", year, month)
	query := `
		SELECT SUM(amount)
		FROM transactions
		WHERE category = ? 
		AND date LIKE ?
		AND amount < 0; 
	`
	var total sql.NullFloat64
	err := db.QueryRow(query, category, dateFilter).Scan(&total)
	if err != nil {
		return 0, err
	}
	if total.Valid {
		return -total.Float64, nil
	}
	return 0, nil
}
