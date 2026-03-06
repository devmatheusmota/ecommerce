package usecase

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	Token    string
	ExpireAt time.Time
}

type LoginUser struct {
	repo repository.UserRepository
}

func NewLoginUser(repo repository.UserRepository) *LoginUser {
	return &LoginUser{repo: repo}
}

func (u *LoginUser) Execute(in LoginUserInput) (*LoginUserOutput, error) {
	user, err := u.repo.GetByEmail(in.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	expireAt := time.Now().Add(time.Hour * 24)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": expireAt.Unix(),
	}).SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	return &LoginUserOutput{Token: token, ExpireAt: expireAt}, nil
}
