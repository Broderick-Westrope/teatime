package db

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/Broderick-Westrope/teatime/internal/entity"
	"github.com/Broderick-Westrope/teatime/internal/secure"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db          *sql.DB
	argonParams *secure.ArgonParams
	keyLength   uint32
}

func NewRepository(dataSourceName string) (*Repository, error) {
	db, err := initDB(dataSourceName)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
		argonParams: &secure.ArgonParams{
			Memory:      64 * 1024,
			Iterations:  1,
			Parallelism: uint8(runtime.NumCPU()),
			SaltLength:  16,
		},
		keyLength: 24,
	}, nil
}

func initDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Create the user_conversations table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS user_conversations (
		username TEXT PRIMARY KEY,
		ciphertext TEXT NOT NULL,
		encryption_params TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("create table failed: %w", err)
	}
	return db, nil
}

func (r *Repository) GetConversations(username, password string) ([]entity.Conversation, error) {
	uc, err := getUserConversations(r.db, username)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return r.setupNewUser(username, password)
		}
		return nil, fmt.Errorf("failed to get user conversations: %w", err)
	}

	key, err := secure.DeriveKey(password, uc.EncryptionParams, r.keyLength)
	if err != nil {
		return nil, fmt.Errorf("failed to derive encryption key: %w", err)
	}

	decodedCiphertext, err := base64.RawStdEncoding.Strict().DecodeString(uc.Ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	plaintextBytes, err := secure.DecryptAESGCM(key, decodedCiphertext)
	if err != nil {
		if errors.Is(err, secure.ErrFailedToDecrypt) {
			// TODO: return a specific error when this failed because of the incorrect key
			return nil, fmt.Errorf("failed to decrypt ciphertext: %w", err)
		}
		return nil, fmt.Errorf("failed to decrypt ciphertext: %w", err)
	}

	var conversations []entity.Conversation
	err = json.Unmarshal(plaintextBytes, &conversations)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal decrypted conversations: %w", err)
	}

	return conversations, nil
}

func (r *Repository) UpdateConversations(username, password string, conversations []entity.Conversation) error {
	uc, err := getUserConversations(r.db, username)
	if err != nil {
		return err
	}

	key, err := secure.DeriveKey(password, uc.EncryptionParams, r.keyLength)
	if err != nil {
		return fmt.Errorf("failed to derive encryption key: %w", err)
	}

	jsonBytes, err := json.Marshal(&conversations)
	if err != nil {
		return fmt.Errorf("failed to unmarshal decrypted conversations: %w", err)
	}

	ciphertextBytes, err := secure.EncryptAESGCM(key, jsonBytes)
	if err != nil {
		return fmt.Errorf("failed to encrypt ciphertext: %w", err)
	}

	encodedCiphertext := base64.RawStdEncoding.EncodeToString(ciphertextBytes)

	uc.UpdatedAt = time.Now()
	uc.Ciphertext = encodedCiphertext

	return updateUserConversations(r.db, uc)
}

func (r *Repository) setupNewUser(username, password string) ([]entity.Conversation, error) {
	// TODO: once user registration & authn is added this data setup should be done there.

	key, hash, err := secure.CreateKey(password, r.argonParams, r.keyLength)
	if err != nil {
		return nil, err
	}

	initialConversations := make([]entity.Conversation, 0)
	jsonBytes, err := json.Marshal(initialConversations)
	if err != nil {
		return nil, err
	}

	ciphertextBytes, err := secure.EncryptAESGCM(key, jsonBytes)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	uc := &UserConversations{
		Username:         username,
		CreatedAt:        now,
		UpdatedAt:        now,
		Ciphertext:       base64.RawStdEncoding.EncodeToString(ciphertextBytes),
		EncryptionParams: hash,
	}

	err = insertUserConversations(r.db, uc)
	return initialConversations, err
}
