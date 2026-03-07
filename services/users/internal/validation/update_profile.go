package validation

import (
	"strings"

	"github.com/ecommerce/services/users/internal/domain"
)

// UpdateProfileInput holds optional profile fields. Only non-empty fields are updated.
type UpdateProfileInput struct {
	Name  string
	Phone string
	CPF   string
}

// ValidateUpdateProfileInput validates the input for PATCH /me.
// At least one of name, phone, or cpf must be present and non-empty.
// Each present field is validated (name/phone non-empty after trim; CPF valid format and check digits).
func ValidateUpdateProfileInput(in *UpdateProfileInput) error {
	name := strings.TrimSpace(in.Name)
	phone := strings.TrimSpace(in.Phone)
	cpf := strings.TrimSpace(in.CPF)

	if name == "" && phone == "" && cpf == "" {
		return domain.ErrValidation("at least one field (name, phone, cpf) is required")
	}

	if name != "" {
		in.Name = name
	}
	if phone != "" {
		in.Phone = phone
	}
	if cpf != "" {
		if !validCPF(cpf) {
			return domain.ErrValidation("invalid cpf")
		}
		in.CPF = cpf
	}
	return nil
}
