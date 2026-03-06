package usecase

import (
	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserInput struct {
	Email    string
	Password string
	Name     string
	Phone    string
	CPF      string
}

type RegisterUser struct {
	repo repository.UserRepository
}

func NewRegisterUser(repo repository.UserRepository) *RegisterUser {
	return &RegisterUser{repo: repo}
}

func (u *RegisterUser) Execute(in RegisterUserInput) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        in.Email,
		Name:         in.Name,
		Phone:        in.Phone,
		CPF:          in.CPF,
		PasswordHash: string(hash),
	}

	return u.repo.Create(user)
}
