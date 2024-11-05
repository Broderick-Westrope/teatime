package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UserConversations struct {
	Username         string
	Ciphertext       string // Base64 encoded, encrypted JSON data containing all conversations
	EncryptionParams string // Encoded string containing parameters for deriving the encryption key from the password

	CreatedAt time.Time
	UpdatedAt time.Time
}

func insertUserConversations(db *sql.DB, uc *UserConversations) error {
	query := `
	INSERT INTO user_conversations (username, ciphertext, encryption_params, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(username) DO UPDATE SET
		updated_at=excluded.updated_at,
		ciphertext=excluded.ciphertext,
		encryption_params=excluded.encryption_params;
	`
	_, err := db.Exec(query, uc.Username, uc.Ciphertext, uc.EncryptionParams, uc.CreatedAt, uc.UpdatedAt)
	return err
}

func getUserConversations(db *sql.DB, username string) (*UserConversations, error) {
	query := `
	SELECT username, ciphertext, encryption_params, created_at, updated_at
	FROM user_conversations
	WHERE username = ?
	`
	row := db.QueryRow(query, username)

	var result UserConversations
	err := row.Scan(&result.Username, &result.Ciphertext, &result.EncryptionParams, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: no user found with username %q: %w", ErrNotFound, username, err)
		}
		return nil, err
	}
	return &result, nil
}

func updateUserConversations(db *sql.DB, uc *UserConversations) error {
	updateSQL := `
	UPDATE user_conversations
	SET ciphertext = ?, encryption_params = ?, updated_at = ?
	WHERE username = ?
	`
	_, err := db.Exec(updateSQL, uc.Ciphertext, uc.EncryptionParams, uc.UpdatedAt, uc.Username)
	return err
}

//nolint:unused // Will be used in the future for erasing user data.
func deleteUserConversations(db *sql.DB, username string) error {
	deleteSQL := `DELETE FROM user_conversations WHERE username = ?`
	_, err := db.Exec(deleteSQL, username)
	return err
}
