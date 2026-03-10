package validation

import (
	"testing"
)

func TestValidateAddressInput(t *testing.T) {
	t.Run("missing street", func(t *testing.T) {
		in := &AddressInput{Number: "1", Neighborhood: "X", City: "Y", State: "SP", ZipCode: "01310100", Type: "shipping"}
		err := ValidateAddressInput(in)
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != "street is required" {
			t.Errorf("got %q", err.Error())
		}
	})
	t.Run("invalid zip_code", func(t *testing.T) {
		in := &AddressInput{Street: "Rua A", Number: "1", Neighborhood: "X", City: "Y", State: "SP", ZipCode: "invalid", Type: "shipping"}
		err := ValidateAddressInput(in)
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != "invalid zip_code format (use 12345-678 or 12345678)" {
			t.Errorf("got %q", err.Error())
		}
	})
	t.Run("invalid type", func(t *testing.T) {
		in := &AddressInput{Street: "Rua A", Number: "1", Neighborhood: "X", City: "Y", State: "SP", ZipCode: "01310100", Type: "other"}
		err := ValidateAddressInput(in)
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != "type must be billing or shipping" {
			t.Errorf("got %q", err.Error())
		}
	})
	t.Run("valid", func(t *testing.T) {
		in := &AddressInput{Street: " Rua Augusta ", Number: " 100 ", Complement: " Sala 1 ", Neighborhood: "Consolação", City: "São Paulo", State: "SP", ZipCode: "01310-100", Type: "shipping"}
		err := ValidateAddressInput(in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if in.Street != "Rua Augusta" || in.ZipCode != "01310-100" || in.Type != "shipping" {
			t.Errorf("got %+v", in)
		}
	})
}
