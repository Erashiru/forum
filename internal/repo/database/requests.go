package database

import (
	"errors"
	"fmt"
	"forum/models"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
)

func (m *Storage) CreateRequest(userid int) error {
	stmt := `INSERT INTO requests (user_id, created_at) VALUES (?, ?)`
	_, err := m.DB.Exec(stmt, userid, time.Now())
	if err != nil {
		fmt.Println(err)
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			if sqliteError.Code == 19 && strings.Contains(sqliteError.Error(), "UNIQUE constraint failed:") {
				return models.ErrDuplicateRequest
			}
		}
		return err
	}
	return nil
}

func (m *Storage) GetRequests() ([]*models.Request, error) {
	stmt := `SELECT user_id, created_at FROM requests`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*models.Request
	for rows.Next() {
		var r models.Request
		err := rows.Scan(&r.SenderID, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		r.SenderName, err = m.GetUser(r.SenderID)
		if err != nil {
			return nil, err
		}

		requests = append(requests, &r)
	}

	return requests, nil
}

func (m *Storage) RequestDone(id int) error {
	stmt := `DELETE FROM requests WHERE user_id = ?`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
