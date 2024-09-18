package database

import (
	"database/sql"
	"errors"
	"forum/models"
)

func (m *Storage) GetRole(id int) (string, error) {
	stmt := `SELECT role FROM roles WHERE user_id = ?`
	rows := m.DB.QueryRow(stmt, id)
	var role string
	err := rows.Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrNoRecord
		} else {
			return "", err
		}
	}
	return role, nil
}

func (m *Storage) CreateUserRole(id int, role string) error {
	stmt := `INSERT INTO roles (user_id, role) VALUES (?. ?)`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *Storage) UpdateUserRole(id int, role string) error {
	stmt := `UPDATE roles SET role = ? WHERE user_id = ?`
	_, err := m.DB.Exec(stmt, role, id)
	if err != nil {
		return err
	}
	err = m.RequestDone(id)
	if err != nil {
		return err
	}
	return nil
}
