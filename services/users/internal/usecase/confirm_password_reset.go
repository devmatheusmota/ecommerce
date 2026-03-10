package usecase

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type ConfirmPasswordResetInput struct {
	Token       string
	NewPassword string
}

type ConfirmPasswordReset struct {
	userRepository          repository.UserRepository
	passwordResetRepository repository.PasswordResetRepository
}

func NewConfirmPasswordReset(userRepository repository.UserRepository, passwordResetRepository repository.PasswordResetRepository) *ConfirmPasswordReset {
	return &ConfirmPasswordReset{
		userRepository:          userRepository,
		passwordResetRepository: passwordResetRepository,
	}
}

func (u *ConfirmPasswordReset) Execute(in ConfirmPasswordResetInput) error {
	hash := sha256.Sum256([]byte(in.Token))
	tokenHash := fmt.Sprintf("%x", hash)
	row, err := u.passwordResetRepository.FindByTokenHash(tokenHash)
	if err != nil {
		return err
	}
	if row == nil {
		return domain.ErrInvalidResetToken
	}
	if time.Now().After(row.ExpiresAt) {
		_ = u.passwordResetRepository.DeleteByID(row.ID)
		return domain.ErrInvalidResetToken
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := u.userRepository.UpdatePassword(row.UserID, string(passwordHash)); err != nil {
		return err
	}
	_ = u.passwordResetRepository.DeleteByID(row.ID)
	return nil
}
