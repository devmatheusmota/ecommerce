package usecase

import (
	"errors"
	"testing"

	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/repository"
)

func TestRegisterUser_Execute(t *testing.T) {
	validCPF := "529.982.247-25"

	t.Run("success", func(t *testing.T) {
		repo := repository.NewMockUserRepository()
		uc := NewRegisterUser(repo)

		user, err := uc.Execute(RegisterUserInput{
			Email:    "user@example.com",
			Password: "password123",
			Name:     "John Doe",
			Phone:    "11999999999",
			CPF:      validCPF,
		})
		if err != nil {
			t.Fatalf("Execute() err = %v", err)
		}
		if user.ID == "" {
			t.Error("expected non-empty ID")
		}
		if user.Email != "user@example.com" || user.Name != "John Doe" || user.Phone != "11999999999" || user.CPF != validCPF {
			t.Errorf("user fields: got %+v", user)
		}
		if user.PasswordHash == "" || user.PasswordHash == "password123" {
			t.Error("password must be hashed")
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		repo := repository.NewMockUserRepository()
		uc := NewRegisterUser(repo)

		input := RegisterUserInput{
			Email:    "same@example.com",
			Password: "password123",
			Name:     "John",
			Phone:    "11999999999",
			CPF:      validCPF,
		}
		_, err := uc.Execute(input)
		if err != nil {
			t.Fatalf("first Execute() err = %v", err)
		}
		_, err = uc.Execute(input)
		if !errors.Is(err, domain.ErrDuplicateEmail) {
			t.Errorf("second Execute() err = %v, want ErrDuplicateEmail", err)
		}
	})
}
