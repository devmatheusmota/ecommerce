package usecase

import (
	"errors"
	"testing"

	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/repository"
)

func TestUpdateProfile_Execute(t *testing.T) {
	validCPF := "529.982.247-25"

	t.Run("success partial update name only", func(t *testing.T) {
		repo := repository.NewMockUserRepository()
		repo.SetUser(&domain.User{
			ID: "user-1", Email: "u@example.com", Name: "Old", Phone: "11999999999",
			CPF: validCPF, PasswordHash: "hash",
		})
		uc := NewUpdateProfile(repo)

		out, err := uc.Execute(UpdateProfileInput{
			UserID: "user-1",
			Name:   "New Name",
			Phone:  "",
			CPF:    "",
		})
		if err != nil {
			t.Fatalf("Execute() err = %v", err)
		}
		if out.Name != "New Name" {
			t.Errorf("name: got %q, want New Name", out.Name)
		}
		if out.Phone != "11999999999" {
			t.Errorf("phone should be unchanged: got %q", out.Phone)
		}
		if out.CPF != validCPF {
			t.Errorf("cpf should be unchanged: got %q", out.CPF)
		}
		if out.Email != "u@example.com" {
			t.Errorf("email: got %q", out.Email)
		}
	})

	t.Run("success update all profile fields", func(t *testing.T) {
		repo := repository.NewMockUserRepository()
		repo.SetUser(&domain.User{
			ID: "user-2", Email: "u2@example.com", Name: "Old", Phone: "11000000000",
			CPF: validCPF, PasswordHash: "hash",
		})
		uc := NewUpdateProfile(repo)

		out, err := uc.Execute(UpdateProfileInput{
			UserID: "user-2",
			Name:   "Jane Doe",
			Phone:  "11888887777",
			CPF:    "100.000.001-08",
		})
		if err != nil {
			t.Fatalf("Execute() err = %v", err)
		}
		if out.Name != "Jane Doe" || out.Phone != "11888887777" || out.CPF != "100.000.001-08" {
			t.Errorf("got name=%q phone=%q cpf=%q", out.Name, out.Phone, out.CPF)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		repo := repository.NewMockUserRepository()
		uc := NewUpdateProfile(repo)

		_, err := uc.Execute(UpdateProfileInput{
			UserID: "nonexistent",
			Name:   "Any",
			Phone:  "",
			CPF:    "",
		})
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})
}
