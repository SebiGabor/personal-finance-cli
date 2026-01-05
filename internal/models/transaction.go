package models

import (
	"database/sql"
	"time"
)

type Transaction struct {
	ID          int64
	Date        time.Time
	Description string
	Amount      float64
	Category    string
	CreatedAt   time.Time
}

// CreateTransaction inserts a new transaction.
func CreateTransaction(db *sql.DB, t *Transaction) error {
	query := `
        INSERT INTO transactions (date, description, amount, category)
        VALUES (?, ?, ?, ?);
    `

	res, err := db.Exec(query,
		t.Date.Format("2006-01-02"),
		t.Description,
		t.Amount,
		t.Category,
	)
	if err != nil {
		return err
	}

	t.ID, err = res.LastInsertId()
	return err
}

// GetTransaction retrieves one by ID
func GetTransaction(db *sql.DB, id int64) (*Transaction, error) {
	query := `
        SELECT id, date, description, amount, category, created_at
        FROM transactions WHERE id = ?;
    `

	row := db.QueryRow(query, id)

	var t Transaction
	var dateStr string

	if err := row.Scan(&t.ID, &dateStr, &t.Description, &t.Amount, &t.Category, &t.CreatedAt); err != nil {
		return nil, err
	}

	t.Date, _ = time.Parse("2006-01-02", dateStr)
	return &t, nil
}

// ListTransactions retrieves all (optionally filtered later)
func ListTransactions(db *sql.DB) ([]Transaction, error) {
	query := `
        SELECT id, date, description, amount, category, created_at
        FROM transactions ORDER BY date DESC;
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Transaction

	for rows.Next() {
		var t Transaction
		var dateStr string

		if err := rows.Scan(&t.ID, &dateStr, &t.Description, &t.Amount, &t.Category, &t.CreatedAt); err != nil {
			return nil, err
		}

		t.Date, _ = time.Parse("2006-01-02", dateStr)
		list = append(list, t)
	}

	return list, nil
}

// UpdateTransaction (simple)
func UpdateTransaction(db *sql.DB, t *Transaction) error {
	query := `
        UPDATE transactions
        SET date = ?, description = ?, amount = ?, category = ?
        WHERE id = ?;
    `

	_, err := db.Exec(query,
		t.Date.Format("2006-01-02"),
		t.Description,
		t.Amount,
		t.Category,
		t.ID,
	)
	return err
}

// DeleteTransaction
func DeleteTransaction(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM transactions WHERE id = ?`, id)
	return err
}

// SearchTransactions filters transactions where description or category matches the query.
func SearchTransactions(db *sql.DB, queryStr string) ([]Transaction, error) {
	// We use the LIKE operator for simple keyword matching
	sqlQuery := `
        SELECT id, date, description, amount, category, created_at
        FROM transactions 
        WHERE description LIKE ? OR category LIKE ?
        ORDER BY date DESC;
    `
	searchTerm := "%" + queryStr + "%"
	rows, err := db.Query(sqlQuery, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Transaction
	for rows.Next() {
		var t Transaction
		var dateStr string
		if err := rows.Scan(&t.ID, &dateStr, &t.Description, &t.Amount, &t.Category, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.Date, _ = time.Parse("2006-01-02", dateStr)
		list = append(list, t)
	}
	return list, nil
}
