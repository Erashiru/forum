package database

import (
	"fmt"
	"forum/models"
)

func (m *Storage) CreateReport(report *models.Report) error {
	stmt := `INSERT INTO reports (user_id, report_text, post_id, created_at)
	VALUES (?, ?, ?, DATETIME('now'))`

	_, err := m.DB.Exec(stmt, report.ModerID, report.Text, report.PostID)
	if err != nil {
		return err
	}

	return nil
}

func (m *Storage) GetReports() ([]*models.Report, error) {
	stmt := `SELECT user_id, report_text, post_id, created_at FROM reports`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*models.Report
	for rows.Next() {
		var r models.Report
		err := rows.Scan(&r.ModerID, &r.Text, &r.PostID, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		r.ModerName, err = m.GetUser(r.ModerID)
		if err != nil {
			return nil, err
		}

		reports = append(reports, &r)
	}

	return reports, nil
}

func (m *Storage) ReportDone(id int) error {
	stmt := `DELETE FROM reports WHERE post_id = ?`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	fmt.Println("aboba")
	return nil
}
