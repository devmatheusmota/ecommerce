package validation

import (
	"strings"

	"github.com/ecommerce/services/users/internal/domain"
)

func ValidateRequestPasswordResetEmail(email string) (string, error) {
	return ValidateEmail(email)
}

func ValidateConfirmPasswordReset(token, newPassword string) error {
	if strings.TrimSpace(token) == "" {
		return domain.ErrValidation("token is required")
	}
	if strings.TrimSpace(newPassword) == "" {
		return domain.ErrValidation("new_password is required")
	}
	if len(newPassword) < 6 {
		return domain.ErrValidation("new_password must be at least 6 characters")
	}
	return nil
}
