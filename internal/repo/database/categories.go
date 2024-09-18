package database

import (
	"encoding/json"
	"fmt"
)

func (m *Storage) ChooseCategories(postid int, categorie []string) error {
	stmt := `INSERT INTO categories (post_id, name)
	VALUES(?, ?);`

	categorieJSON, err := json.Marshal(categorie)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(stmt, postid, string(categorieJSON))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (m *Storage) GetCategory(postid int) ([]string, error) {
	stmt := `SELECT name FROM categories WHERE post_id = ?`

	row := m.DB.QueryRow(stmt, postid)
	var category string
	err := row.Scan(&category)
	if err != nil {
		return nil, err
	}
	var cata []string

	err = json.Unmarshal([]byte(category), &cata)
	if err != nil {
		return nil, err
	}
	return cata, nil
}
