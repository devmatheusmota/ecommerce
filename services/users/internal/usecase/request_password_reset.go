package usecase

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/repository"
)

const resetTokenExpiry = time.Hour
const resetTokenBytes = 32

type RequestPasswordResetInput struct {
	Email string
}

type RequestPasswordResetOutput struct {
	Token     string // For dev/testing; in prod this would be sent by email only
	ExpiresAt time.Time
}

type RequestPasswordReset struct {
	userRepository          repository.UserRepository
	passwordResetRepository repository.PasswordResetRepository
}

func NewRequestPasswordReset(userRepository repository.UserRepository, passwordResetRepository repository.PasswordResetRepository) *RequestPasswordReset {
	return &RequestPasswordReset{
		userRepository:          userRepository,
		passwordResetRepository: passwordResetRepository,
	}
}

func (u *RequestPasswordReset) Execute(in RequestPasswordResetInput) (*RequestPasswordResetOutput, error) {
	user, err := u.userRepository.GetByEmail(in.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	tokenBytes := make([]byte, resetTokenBytes)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	hash := sha256.Sum256([]byte(token))
	tokenHash := fmt.Sprintf("%x", hash)
	expiresAt := time.Now().Add(resetTokenExpiry)
	if err := u.passwordResetRepository.Create(user.ID, tokenHash, expiresAt); err != nil {
		return nil, err
	}
	return &RequestPasswordResetOutput{Token: token, ExpiresAt: expiresAt}, nil
}
