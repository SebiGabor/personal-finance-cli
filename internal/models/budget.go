package models

import (
	"database/sql"
)

type Budget struct {
	ID       int64
	Category string
	Amount   float64
	Period   string // "monthly", "weekly", "yearly"
}

func CreateBudget(db *sql.DB, b *Budget) error {
	query := `
        INSERT INTO budgets (category, amount, period)
        VALUES (?, ?, ?);
    `
	res, err := db.Exec(query, b.Category, b.Amount, b.Period)
	if err != nil {
		return err
	}
	b.ID, err = res.LastInsertId()
	return err
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
	rows, err := db.Query(`
        SELECT id, category, amount, period
        FROM budgets ORDER BY category;
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Budget
	for rows.Next() {
		var b Budget
		if err := rows.Scan(&b.ID, &b.Category, &b.Amount, &b.Period); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
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
	_, err := db.Exec(`DELETE FROM budgets WHERE id = ?`, id)
	return err
}
