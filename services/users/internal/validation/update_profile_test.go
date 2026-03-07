package validation

import (
	"errors"
	"testing"

	"github.com/ecommerce/services/users/internal/domain"
)

func TestValidateUpdateProfileInput(t *testing.T) {
	validCPF := "529.982.247-25"

	t.Run("at least one field required", func(t *testing.T) {
		in := &UpdateProfileInput{Name: "", Phone: "", CPF: ""}
		err := ValidateUpdateProfileInput(in)
		if err == nil {
			t.Fatal("expected validation error")
		}
		var validationErr domain.ErrValidation
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ErrValidation, got %T", err)
		}
		if err.Error() != "at least one field (name, phone, cpf) is required" {
			t.Errorf("error: got %q", err.Error())
		}
	})

	t.Run("invalid cpf", func(t *testing.T) {
		in := &UpdateProfileInput{Name: "", Phone: "", CPF: "111.111.111-11"}
		err := ValidateUpdateProfileInput(in)
		if err == nil {
			t.Fatal("expected validation error for invalid CPF")
		}
		if err.Error() != "invalid cpf" {
			t.Errorf("error: got %q", err.Error())
		}
	})

	t.Run("valid name only", func(t *testing.T) {
		in := &UpdateProfileInput{Name: "  New Name  ", Phone: "", CPF: ""}
		err := ValidateUpdateProfileInput(in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if in.Name != "New Name" {
			t.Errorf("name: got %q, want trimmed", in.Name)
		}
	})

	t.Run("valid phone only", func(t *testing.T) {
		in := &UpdateProfileInput{Name: "", Phone: " 11999999999 ", CPF: ""}
		err := ValidateUpdateProfileInput(in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if in.Phone != "11999999999" {
			t.Errorf("phone: got %q", in.Phone)
		}
	})

	t.Run("valid cpf only", func(t *testing.T) {
		in := &UpdateProfileInput{Name: "", Phone: "", CPF: validCPF}
		err := ValidateUpdateProfileInput(in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if in.CPF != validCPF {
			t.Errorf("cpf: got %q", in.CPF)
		}
	})

	t.Run("valid all fields", func(t *testing.T) {
		in := &UpdateProfileInput{Name: "Jane", Phone: "11888887777", CPF: "100.000.001-08"}
		err := ValidateUpdateProfileInput(in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if in.Name != "Jane" || in.Phone != "11888887777" || in.CPF != "100.000.001-08" {
			t.Errorf("got %+v", in)
		}
	})
}
