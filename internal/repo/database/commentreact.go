package database

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/models"
)

func (m *Storage) CreateCommentReaction(userid, commentid, reaction int) error {
	stmt := `SELECT reaction_status FROM comment_reactions WHERE user_id = ? AND comment_id = ?`
	var is int
	err := m.DB.QueryRow(stmt, userid, commentid).Scan(&is)
	if err != nil {
		if err == sql.ErrNoRows {
			stmt = `INSERT INTO comment_reactions (user_id, comment_id, reaction_status) 
			VALUES (?, ?, ?)`

			_, err = m.DB.Exec(stmt, userid, commentid, reaction)
			if err != nil {
				fmt.Println(err)
				return err
			}

		} else {
			return err
		}
	} else {
		if is == reaction {
			stmt = `DELETE FROM comment_reactions WHERE user_id = ? AND comment_id = ?`
			_, err := m.DB.Exec(stmt, userid, commentid)
			if err != nil {
				return err
			}
		} else {
			stmt = `UPDATE comment_reactions SET reaction_status = ? WHERE user_id = ? AND comment_id = ?`
			_, err := m.DB.Exec(stmt, reaction, userid, commentid)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Storage) GetCommentLikes(commentid int) (int, error) {
	stmt := `SELECT COUNT(reaction_status) FROM comment_reactions WHERE reaction_status = 1 AND comment_id = ?`

	row := m.DB.QueryRow(stmt, commentid)

	var num int

	err := row.Scan(&num)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		} else {
			return 0, err
		}
	}
	return num, nil
}

func (m *Storage) GetCommentDislikes(commentid int) (int, error) {
	stmt := `SELECT COUNT(reaction_status) FROM comment_reactions WHERE reaction_status = -1 AND comment_id = ?`

	row := m.DB.QueryRow(stmt, commentid)

	var num int

	err := row.Scan(&num)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		} else {
			return 0, err
		}
	}
	return num, nil
}
