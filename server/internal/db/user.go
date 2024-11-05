package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	Username     string
	PasswordHash string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func insertUser(db *sql.DB, user *User) error {
	query := `
	INSERT INTO users (username, password_hash, created_at, updated_at)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT(username) DO UPDATE SET
		password_hash=excluded.password_hash,
		updated_at=excluded.updated_at;
	`
	_, err := db.Exec(query, user.Username, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	return err
}

func getUser(db *sql.DB, username string) (*User, error) {
	query := `
	SELECT username, password_hash, created_at, updated_at
	FROM users
	WHERE username = $1
	`
	row := db.QueryRow(query, username)

	var result User
	err := row.Scan(&result.Username, &result.PasswordHash, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: no user found with username %q: %w", ErrNotFound, username, err)
		}
		return nil, err
	}
	return &result, nil
}

func updateUser(db *sql.DB, user *User) error {
	updateSQL := `
	UPDATE users
	SET password_hash = $1, updated_at = $2
	WHERE username = $3
	`
	_, err := db.Exec(updateSQL, user.PasswordHash, user.UpdatedAt, user.Username)
	return err
}

func deleteUser(db *sql.DB, username string) error {
	deleteSQL := `DELETE FROM users WHERE username = $1`
	_, err := db.Exec(deleteSQL, username)
	return err
}
