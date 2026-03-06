package validation

import (
	"strings"

	"github.com/ecommerce/services/users/internal/domain"
)

// ValidateLoginInput checks email (required, format) and password (required).
// Returns the normalized email (trim, lowercase) for the handler to pass to the use case.
func ValidateLoginInput(email, password string) (normalizedEmail string, err error) {
	normalized, err := ValidateEmail(email)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(password) == "" {
		return "", domain.ErrValidation("password is required")
	}
	return normalized, nil
}
