package repository

import (
	"database/sql"
	"errors"
	"time"
)

type PostgresPasswordResetRepository struct {
	db *sql.DB
}

func NewPostgresPasswordResetRepository(db *sql.DB) *PostgresPasswordResetRepository {
	return &PostgresPasswordResetRepository{db: db}
}

func (r *PostgresPasswordResetRepository) Create(userID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(`
		INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt)
	return err
}

func (r *PostgresPasswordResetRepository) FindByTokenHash(tokenHash string) (*PasswordResetToken, error) {
	var t PasswordResetToken
	err := r.db.QueryRow(`
		SELECT id, user_id, token_hash, expires_at
		FROM password_reset_tokens
		WHERE token_hash = $1
	`, tokenHash).Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *PostgresPasswordResetRepository) DeleteByID(id string) error {
	_, err := r.db.Exec(`DELETE FROM password_reset_tokens WHERE id = $1`, id)
	return err
}
