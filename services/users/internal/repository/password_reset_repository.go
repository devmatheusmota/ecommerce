package repository

import "time"

type PasswordResetToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetRepository interface {
	Create(userID, tokenHash string, expiresAt time.Time) error
	FindByTokenHash(tokenHash string) (*PasswordResetToken, error)
	DeleteByID(id string) error
}
