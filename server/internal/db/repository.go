package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/redis/go-redis/v9"

	"github.com/Broderick-Westrope/teatime/internal/secure"

	// Postgres database driver.
	_ "github.com/lib/pq"
)

type Repository struct {
	db          *sql.DB
	argonParams *argon2id.Params

	redis              *redis.Client
	redisUserPrefix    string
	redisSessionPrefix string
}

func NewRepository(dbConn, redisAddr string) (*Repository, error) {
	db, err := initDB(dbConn)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
		argonParams: &argon2id.Params{
			Memory:      64 * 1024,
			Iterations:  3,
			Parallelism: 1,
			SaltLength:  16,
			KeyLength:   32,
		},
		redis: redis.NewClient(&redis.Options{
			Addr: redisAddr,
		}),
		redisUserPrefix:    "user:",
		redisSessionPrefix: "session:",
	}, nil
}

func initDB(dbConn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return nil, err
	}

	// Create the user_conversations table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);
	`
	if _, err = db.Exec(createTableSQL); err != nil {
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

func (r *Repository) GetNewSessionID(ctx context.Context, username string) (string, error) {
	oldSessionID, err := r.redis.Get(ctx, r.redisUserPrefix+username).Result()
	if err == nil {
		// Delete the old session ID
		err = r.redis.Del(ctx, r.redisSessionPrefix+oldSessionID).Err()
		if err != nil {
			return "", err
		}
	}

	sessionID, err := secure.GenerateSessionID()
	if err != nil {
		return "", err
	}

	sessionExpiration := time.Hour * 3
	err = r.redis.Set(ctx, r.redisSessionPrefix+sessionID, username, sessionExpiration).Err()
	if err != nil {
		return "", err
	}

	err = r.redis.Set(ctx, r.redisUserPrefix+username, sessionID, sessionExpiration).Err()
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (r *Repository) GetUsernameWithSessionID(ctx context.Context, sessionID string) (string, error) {
	username, err := r.redis.Get(ctx, r.redisSessionPrefix+sessionID).Result()
	if err != nil {
		return "", err
	}

	storedSessionID, err := r.redis.Get(ctx, r.redisUserPrefix+username).Result()
	switch {
	case err != nil:
		return "", err
	case storedSessionID != sessionID:
		return "", errors.New("provided session ID does not match stored session ID")
	}
	return username, nil
}

func (r *Repository) DeleteUserSessions(ctx context.Context, username string) error {
	sessionID, err := r.redis.Get(ctx, r.redisUserPrefix+username).Result()
	if err != nil {
		return err
	}

	err = r.redis.Del(ctx, r.redisUserPrefix+username).Err()
	if err != nil {
		return err
	}
	err = r.redis.Del(ctx, r.redisSessionPrefix+sessionID).Err()
	if err != nil {
		return err
	}
	return nil
}
