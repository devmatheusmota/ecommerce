package validation

import (
	"strings"
	"testing"
)

func TestValidateRegisterInput(t *testing.T) {
	validCPF := "529.982.247-25" // valid check digits

	tests := []struct {
		name        string
		in          RegisterInput
		wantErr     bool
		errContains string
	}{
		{
			name:        "missing email",
			in:          RegisterInput{Password: "secret12", Name: "John", Phone: "11999999999", CPF: validCPF},
			wantErr:     true,
			errContains: "email",
		},
		{
			name:        "invalid email",
			in:          RegisterInput{Email: "bad", Password: "secret12", Name: "John", Phone: "11999999999",  CPF: validCPF},
			wantErr:     true,
			errContains: "email",
		},
		{
			name:        "short password",
			in:          RegisterInput{Email: "a@b.com", Password: "short", Name: "John", Phone: "11999999999",  CPF: validCPF},
			wantErr:     true,
			errContains: "password",
		},
		{
			name:        "empty name",
			in:          RegisterInput{Email: "a@b.com", Password: "secret12", Name: "  ", Phone: "11999999999",  CPF: validCPF},
			wantErr:     true,
			errContains: "name",
		},
		{
			name:        "empty phone",
			in:          RegisterInput{Email: "a@b.com", Password: "secret12", Name: "John", Phone: "  ",  CPF: validCPF},
			wantErr:     true,
			errContains: "phone",
		},
		{
			name:        "empty cpf",
			in:          RegisterInput{Email: "a@b.com", Password: "secret12", Name: "John", Phone: "11999999999",  CPF: "  "},
			wantErr:     true,
			errContains: "cpf",
		},
		{
			name:        "invalid cpf all same digit",
			in:          RegisterInput{Email: "a@b.com", Password: "secret12", Name: "John", Phone: "11999999999",  CPF: "111.111.111-11"},
			wantErr:     true,
			errContains: "cpf",
		},
		{
			name:        "invalid cpf wrong length",
			in:          RegisterInput{Email: "a@b.com", Password: "secret12", Name: "John", Phone: "11999999999",  CPF: "111.111.111-1"},
			wantErr:     true,
			errContains: "cpf",
		},
		{
			name:        "invalid cpf wrong first check digit",
			in:          RegisterInput{Email: "a@b.com", Password: "secret12", Name: "John", Phone: "11999999999",  CPF: "529.982.247-15"},
			wantErr:     true,
			errContains: "cpf",
		},
		{
			name:        "invalid cpf wrong second check digit",
			in:          RegisterInput{Email: "a@b.com", Password: "secret12", Name: "John", Phone: "11999999999",  CPF: "529.982.247-24"},
			wantErr:     true,
			errContains: "cpf",
		},
		{
			name:    "valid",
			in:      RegisterInput{Email: "  user@Example.COM  ", Password: "secret12", Name: "John Doe", Phone: "11999999999",  CPF: validCPF},
			wantErr: false,
		},
		{
			name:    "valid cpf with first check digit 0",
			in:      RegisterInput{Email: "a@b.com", Password: "secret12", Name: "John", Phone: "11999999999",  CPF: "100.000.001-08"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRegisterInput(&tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRegisterInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("ValidateRegisterInput() err = %q, want containing %q", err.Error(), tt.errContains)
			}
			if !tt.wantErr && tt.name == "valid" && tt.in.Email != "user@example.com" {
				t.Errorf("expected email normalized to user@example.com, got %q", tt.in.Email)
			}
		})
	}
}
