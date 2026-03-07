package usecase

import (
	"errors"
	"testing"

	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func TestMeUser_Execute(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:           "user-123",
		Email:        "me@example.com",
		Name:         "Me User",
		Phone:        "11988887777",
		CPF:          "529.982.247-25",
		PasswordHash: string(hash),
	}

	t.Run("success", func(t *testing.T) {
		repo := repository.NewMockUserRepository()
		repo.SetUser(existingUser)
		uc := NewMeUser(repo)

		out, err := uc.Execute(MeUserInput{UserID: "user-123"})
		if err != nil {
			t.Fatalf("Execute() err = %v", err)
		}
		if out.ID != "user-123" || out.Email != "me@example.com" || out.Name != "Me User" || out.Phone != "11988887777" {
			t.Errorf("got %+v", out)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := repository.NewMockUserRepository()
		uc := NewMeUser(repo)

		_, err := uc.Execute(MeUserInput{UserID: "nonexistent"})
		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Errorf("Execute() err = %v, want ErrUserNotFound", err)
		}
	})
}
