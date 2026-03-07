package validation

import (
	"strings"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		wantEmail   string
		wantErr     bool
		errContains string
	}{
		{"empty", "", "", true, "email is required"},
		{"whitespace", "  ", "", true, "email is required"},
		{"invalid format", "notanemail", "", true, "invalid email format"},
		{"missing at", "user.domain.com", "", true, "invalid email format"},
		{"valid lowercase", "user@example.com", "user@example.com", false, ""},
		{"valid normalized", "  User@Example.COM  ", "user@example.com", false, ""},
		{"valid with plus", "user+tag@example.com", "user+tag@example.com", false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContains != "" && (err == nil || !strings.Contains(err.Error(), tt.errContains)) {
				if err != nil {
					t.Errorf("ValidateEmail() err = %q, want containing %q", err.Error(), tt.errContains)
				}
				return
			}
			if !tt.wantErr && got != tt.wantEmail {
				t.Errorf("ValidateEmail() = %q, want %q", got, tt.wantEmail)
			}
		})
	}
}
