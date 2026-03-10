package validation

import (
	"testing"
)

func TestValidateRequestPasswordResetEmail(t *testing.T) {
	t.Run("empty email", func(t *testing.T) {
		_, err := ValidateRequestPasswordResetEmail("")
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != "email is required" {
			t.Errorf("got %q", err.Error())
		}
	})
	t.Run("invalid email", func(t *testing.T) {
		_, err := ValidateRequestPasswordResetEmail("not-an-email")
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != "invalid email format" {
			t.Errorf("got %q", err.Error())
		}
	})
	t.Run("valid email", func(t *testing.T) {
		email, err := ValidateRequestPasswordResetEmail("user@example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if email != "user@example.com" {
			t.Errorf("got %q", email)
		}
	})
}

func TestValidateConfirmPasswordReset(t *testing.T) {
	t.Run("empty token", func(t *testing.T) {
		err := ValidateConfirmPasswordReset("", "newpass123")
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != "token is required" {
			t.Errorf("got %q", err.Error())
		}
	})
	t.Run("empty new_password", func(t *testing.T) {
		err := ValidateConfirmPasswordReset("abc123", "")
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != "new_password is required" {
			t.Errorf("got %q", err.Error())
		}
	})
	t.Run("short new_password", func(t *testing.T) {
		err := ValidateConfirmPasswordReset("abc123", "12345")
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != "new_password must be at least 6 characters" {
			t.Errorf("got %q", err.Error())
		}
	})
	t.Run("valid", func(t *testing.T) {
		err := ValidateConfirmPasswordReset("token123", "newpass123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
