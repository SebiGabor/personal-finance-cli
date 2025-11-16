package models

import (
	"database/sql"
)

type CategoryRule struct {
	ID       int64
	Pattern  string // regex
	Category string
}

func CreateRule(db *sql.DB, r *CategoryRule) error {
	res, err := db.Exec(`
        INSERT INTO category_rules (pattern, category)
        VALUES (?, ?);
    `, r.Pattern, r.Category)
	if err != nil {
		return err
	}
	r.ID, err = res.LastInsertId()
	return err
}

func ListRules(db *sql.DB) ([]CategoryRule, error) {
	rows, err := db.Query(`
        SELECT id, pattern, category
        FROM category_rules;
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []CategoryRule
	for rows.Next() {
		var r CategoryRule
		if err := rows.Scan(&r.ID, &r.Pattern, &r.Category); err != nil {
			return nil, err
		}
		list = append(list, r)
	}
	return list, nil
}

func DeleteRule(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM category_rules WHERE id = ?`, id)
	return err
}
