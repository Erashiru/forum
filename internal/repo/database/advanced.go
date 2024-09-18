package database

func (m *Storage) DeleteUserPost(postid int) error {
	stmt := `DELETE FROM posts WHERE post_id = ?`

	_, err := m.DB.Exec(stmt, postid)
	if err != nil {
		return err
	}

	stmt = `DELETE FROM reactions WHERE post_id = ?`
	_, err = m.DB.Exec(stmt, postid)
	if err != nil {
		return err
	}

	return nil
}

func (m *Storage) UpdateUserPost(postid int, text string, title string) error {
	stmt := `UPDATE posts SET content = ?, title = ? WHERE post_id = ?`

	_, err := m.DB.Exec(stmt, text, title, postid)
	if err != nil {
		return err
	}

	return nil
}
