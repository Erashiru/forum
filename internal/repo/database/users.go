package database

import (
	"database/sql"
	"errors"
	"forum/models"
	"strings"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// User auth by google
func (m *Storage) SaveUser(form *models.UserSignupForm) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), 12)
	if err != nil {
		return 0, err
	}
	stmt := `
		INSERT INTO users (username, email, hash_password)
		VALUES(?, ?, ?)
	`

	res, err := m.DB.Exec(stmt, form.Name, form.Email, string(hashedPassword))
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			if sqliteError.ExtendedCode == sqlite3.ErrConstraintUnique && strings.Contains(sqliteError.Error(), "users.email") {
				return 0, models.ErrDuplicateEmail
			} else if sqliteError.ExtendedCode == sqlite3.ErrConstraintUnique && strings.Contains(sqliteError.Error(), "users.username") {
				return 0, models.ErrDuplicateName
			}
		}
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *Storage) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	stmt := `SELECT user_id, username FROM users WHERE email = ?`

	err := m.DB.QueryRow(stmt, email).Scan(&user.ID, &user.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrInvalidCredentials
		} else {
			return nil, err
		}
	}
	return &user, nil
}

// default user creation
func (m *Storage) CreateUser(username, email, password string) (int, error) {
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	stmt := `
		INSERT INTO users (username, email, hash_password)
		VALUES(?, ?, ?)
	`

	result, err := m.DB.Exec(stmt, username, email, string(hashedpassword))
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			if sqliteError.ExtendedCode == sqlite3.ErrConstraintUnique && strings.Contains(sqliteError.Error(), "users.email") {
				return 0, models.ErrDuplicateEmail
			} else if sqliteError.ExtendedCode == sqlite3.ErrConstraintUnique && strings.Contains(sqliteError.Error(), "users.username") {
				return 0, models.ErrDuplicateName
			}
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = m.CreateUserRole(int(id), "user")
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Authenticate user by email and password
func (m *Storage) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT user_id, hash_password FROM users WHERE email = ?"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *Storage) Exitsts(name string) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE username = ?)"

	err := m.DB.QueryRow(stmt, name).Scan(&exists)

	return exists, err
}

// Get user by id
func (m *Storage) GetUser(id int) (string, error) {
	var username string
	stmt := "SELECT username FROM users WHERE user_id = ?"
	err := m.DB.QueryRow(stmt, id).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
