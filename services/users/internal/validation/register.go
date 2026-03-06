package validation

import (
	"regexp"
	"strings"

	"github.com/ecommerce/services/users/internal/domain"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type RegisterInput struct {
	Email    string
	Password string
	Name     string
	Phone    string
	CPF      string
}

func ValidateRegisterInput(in *RegisterInput) error {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	if in.Email == "" {
		return domain.ErrValidation("email is required")
	}
	if !emailRegex.MatchString(in.Email) {
		return domain.ErrValidation("invalid email format")
	}
	if len(in.Password) < 6 {
		return domain.ErrValidation("password must be at least 6 characters")
	}
	if strings.TrimSpace(in.Name) == "" {
		return domain.ErrValidation("name is required")
	}
	if strings.TrimSpace(in.Phone) == "" {
		return domain.ErrValidation("phone is required")
	}
	if strings.TrimSpace(in.CPF) == "" {
		return domain.ErrValidation("cpf is required")
	}
	if !validCPF(in.CPF) {
		return domain.ErrValidation("invalid cpf")
	}
	return nil
}
