package validation

// Common holds shared validation logic used by multiple validators (register, login, etc.).
// When adding new validation rules that could apply to more than one flow, add them here
// and reuse in the specific validators instead of duplicating. See doc.go for package policy.

import (
	"regexp"
	"strings"

	"github.com/ecommerce/services/users/internal/domain"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail normalizes the email (trim, lowercase), checks it is non-empty and valid format.
// Returns the normalized email and nil, or ("", error) on validation failure.
// Use this from register, login, or any other validator that needs email validation.
func ValidateEmail(email string) (string, error) {
	normalized := strings.TrimSpace(strings.ToLower(email))
	if normalized == "" {
		return "", domain.ErrValidation("email is required")
	}
	if !emailRegex.MatchString(normalized) {
		return "", domain.ErrValidation("invalid email format")
	}
	return normalized, nil
}
