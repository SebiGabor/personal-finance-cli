package models

import (
	"database/sql"
	"fmt"
)

// CategoryTotal holds the sum of amounts for a specific category
type CategoryTotal struct {
	Category string
	Amount   float64
}

// GetMonthlyReport returns the category breakdown, total income, and total expense for a given month/year.
func GetMonthlyReport(db *sql.DB, year int, month int) ([]CategoryTotal, float64, float64, error) {
	// SQLite stores dates as strings "YYYY-MM-DD", so we filter by the "YYYY-MM" prefix
	dateFilter := fmt.Sprintf("%04d-%02d", year, month)

	query := `
		SELECT category, SUM(amount)
		FROM transactions
		WHERE strftime('%Y-%m', date) = ?
		GROUP BY category
		ORDER BY SUM(amount) ASC;
	`

	rows, err := db.Query(query, dateFilter)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var breakdown []CategoryTotal
	var totalIncome, totalExpense float64

	for rows.Next() {
		var ct CategoryTotal
		if err := rows.Scan(&ct.Category, &ct.Amount); err != nil {
			return nil, 0, 0, err
		}
		breakdown = append(breakdown, ct)

		if ct.Amount > 0 {
			totalIncome += ct.Amount
		} else {
			totalExpense += ct.Amount // This will be negative
		}
	}

	return breakdown, totalIncome, totalExpense, nil
}
