package usecase

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// failingGetByEmailRepo implements UserRepository and returns a fixed error from GetByEmail (for coverage).
type failingGetByEmailRepo struct {
	err error
}

func (f *failingGetByEmailRepo) Create(*domain.User) (*domain.User, error) {
	return nil, errors.New("not implemented")
}
func (f *failingGetByEmailRepo) GetByEmail(string) (*domain.User, error) {
	return nil, f.err
}
func (f *failingGetByEmailRepo) GetByID(string) (*domain.User, error) {
	return nil, errors.New("not implemented")
}
func (f *failingGetByEmailRepo) Update(*domain.User) (*domain.User, error) {
	return nil, errors.New("not implemented")
}
func (f *failingGetByEmailRepo) UpdatePassword(string, string) error {
	return errors.New("not implemented")
}

func TestLoginUser_Execute(t *testing.T) {
	const testJWTSecret = "test-secret-for-login"
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", testJWTSecret)
	t.Cleanup(func() {
		os.Setenv("JWT_SECRET", originalSecret)
	})

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:           "user-456",
		Email:        "login@example.com",
		Name:         "Login User",
		Phone:        "11977776666",
		CPF:          "529.982.247-25",
		PasswordHash: string(hash),
	}

	t.Run("success", func(t *testing.T) {
		repo := repository.NewMockUserRepository()
		repo.SetUser(existingUser)
		uc := NewLoginUser(repo)

		out, err := uc.Execute(LoginUserInput{Email: "login@example.com", Password: "correct"})
		if err != nil {
			t.Fatalf("Execute() err = %v", err)
		}
		if out.Token == "" {
			t.Error("expected non-empty token")
		}
		if out.ExpireAt.Before(time.Now()) {
			t.Error("expire_at should be in the future")
		}
	})

	t.Run("user not found", func(t *testing.T) {
		repo := repository.NewMockUserRepository()
		uc := NewLoginUser(repo)

		_, err := uc.Execute(LoginUserInput{Email: "nobody@example.com", Password: "any"})
		if !errors.Is(err, domain.ErrInvalidCredentials) {
			t.Errorf("Execute() err = %v, want ErrInvalidCredentials", err)
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		repo := repository.NewMockUserRepository()
		repo.SetUser(existingUser)
		uc := NewLoginUser(repo)

		_, err := uc.Execute(LoginUserInput{Email: "login@example.com", Password: "wrong"})
		if !errors.Is(err, domain.ErrInvalidCredentials) {
			t.Errorf("Execute() err = %v, want ErrInvalidCredentials", err)
		}
	})

	t.Run("JWT_SECRET unset", func(t *testing.T) {
		os.Unsetenv("JWT_SECRET")
		t.Cleanup(func() { os.Setenv("JWT_SECRET", testJWTSecret) })

		repo := repository.NewMockUserRepository()
		repo.SetUser(existingUser)
		uc := NewLoginUser(repo)

		_, err := uc.Execute(LoginUserInput{Email: "login@example.com", Password: "correct"})
		if err == nil {
			t.Fatal("Execute() expected error when JWT_SECRET is unset")
		}
		if err != nil && err.Error() != "JWT_SECRET is not set" {
			t.Errorf("Execute() err = %v, want JWT_SECRET is not set", err)
		}
	})

	t.Run("repo GetByEmail returns non-ErrUserNotFound error", func(t *testing.T) {
		dbErr := errors.New("database error")
		uc := NewLoginUser(&failingGetByEmailRepo{err: dbErr})

		_, err := uc.Execute(LoginUserInput{Email: "any@example.com", Password: "any"})
		if !errors.Is(err, dbErr) {
			t.Errorf("Execute() err = %v, want %v", err, dbErr)
		}
	})
}
