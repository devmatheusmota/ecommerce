package validation

import (
	"strings"
	"testing"
)

func TestValidateLoginInput(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		wantEmail   string
		wantErr     bool
		errContains string
	}{
		{"empty email", "", "password", "", true, "email"},
		{"invalid email", "bad", "password", "", true, "email"},
		{"empty password", "user@example.com", "", "", true, "password"},
		{"whitespace password", "user@example.com", "  ", "", true, "password"},
		{"valid", "  User@Example.COM  ", "secret", "user@example.com", false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEmail, err := ValidateLoginInput(tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLoginInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("ValidateLoginInput() err = %q, want containing %q", err.Error(), tt.errContains)
			}
			if !tt.wantErr && gotEmail != tt.wantEmail {
				t.Errorf("ValidateLoginInput() email = %q, want %q", gotEmail, tt.wantEmail)
			}
		})
	}
}
