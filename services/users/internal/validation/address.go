package validation

import (
	"regexp"
	"strings"

	"github.com/ecommerce/services/users/internal/domain"
)

var zipCodeRegex = regexp.MustCompile(`^\d{5}-?\d{3}$`) // Brazilian CEP: 12345-678 or 12345678

type AddressInput struct {
	Street            string
	Number            string
	Complement        string
	Neighborhood      string
	City              string
	State             string
	ZipCode           string
	Type              string
	IsDefaultBilling  bool
	IsDefaultShipping bool
}

func ValidateAddressInput(in *AddressInput) error {
	if strings.TrimSpace(in.Street) == "" {
		return domain.ErrValidation("street is required")
	}
	if strings.TrimSpace(in.Number) == "" {
		return domain.ErrValidation("number is required")
	}
	if strings.TrimSpace(in.Neighborhood) == "" {
		return domain.ErrValidation("neighborhood is required")
	}
	if strings.TrimSpace(in.City) == "" {
		return domain.ErrValidation("city is required")
	}
	if strings.TrimSpace(in.State) == "" {
		return domain.ErrValidation("state is required")
	}
	zip := strings.TrimSpace(strings.ReplaceAll(in.ZipCode, " ", ""))
	if zip == "" {
		return domain.ErrValidation("zip_code is required")
	}
	if !zipCodeRegex.MatchString(zip) {
		return domain.ErrValidation("invalid zip_code format (use 12345-678 or 12345678)")
	}
	in.ZipCode = zip
	typeNorm := strings.TrimSpace(strings.ToLower(in.Type))
	if typeNorm != "billing" && typeNorm != "shipping" {
		return domain.ErrValidation("type must be billing or shipping")
	}
	in.Type = typeNorm
	in.Complement = strings.TrimSpace(in.Complement)
	in.Street = strings.TrimSpace(in.Street)
	in.Number = strings.TrimSpace(in.Number)
	in.Neighborhood = strings.TrimSpace(in.Neighborhood)
	in.City = strings.TrimSpace(in.City)
	in.State = strings.TrimSpace(in.State)
	return nil
}
