package validation

import (
	"strings"

	"github.com/ecommerce/services/users/internal/domain"
)

type RegisterInput struct {
	Email    string
	Password string
	Name     string
	Phone    string
	CPF      string
}

func ValidateRegisterInput(in *RegisterInput) error {
	email, err := ValidateEmail(in.Email)
	if err != nil {
		return err
	}
	in.Email = email
	if len(in.Password) < 6 {
		return domain.ErrValidation("password must be at least 6 characters")
	}
	if strings.TrimSpace(in.Name) == "" {
		return domain.ErrValidation("name is required")
	}
	phone := strings.TrimSpace(in.Phone)
	if phone == "" {
		return domain.ErrValidation("phone is required")
	}
	in.Phone = phone
	if strings.TrimSpace(in.CPF) == "" {
		return domain.ErrValidation("cpf is required")
	}
	if !validCPF(in.CPF) {
		return domain.ErrValidation("invalid cpf")
	}
	return nil
}
