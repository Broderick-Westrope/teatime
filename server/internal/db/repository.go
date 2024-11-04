package db

import (
	"database/sql"
	"fmt"
	"runtime"
	"time"

	"github.com/alexedwards/argon2id"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db          *sql.DB
	argonParams *argon2id.Params
}

func NewRepository(dataSourceName string) (*Repository, error) {
	db, err := initDB(dataSourceName)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
		// TODO: revisit these parameters. currently using defaults
		argonParams: &argon2id.Params{
			Memory:      64 * 1024,
			Iterations:  1,
			Parallelism: uint8(runtime.NumCPU()),
			SaltLength:  16,
			KeyLength:   32,
		},
	}, nil
}

func initDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Create the user_conversations table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("create table failed: %w", err)
	}
	return db, nil
}

func (r *Repository) CreateUser(username, password string) error {
	passwordHash, err := argon2id.CreateHash(password, r.argonParams)
	if err != nil {
		return err
	}

	now := time.Now()
	return insertUser(r.db, &User{
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
}

func (r *Repository) AuthenticateUser(username, password string) (bool, error) {
	user, err := getUser(r.db, username)
	if err != nil {
		return false, err
	}

	match, _, err := argon2id.CheckHash(password, user.PasswordHash)
	return match, err
}
