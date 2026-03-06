package usecase

import (
	"github.com/ecommerce/services/users/internal/repository"
)

type MeUser struct {
	repo repository.UserRepository
}

type MeUserInput struct {
	UserID string
}

type MeUserOutput struct {
	ID    string
	Email string
	Name  string
	Phone string
	CPF   string
}

func NewMeUser(repo repository.UserRepository) *MeUser {
	return &MeUser{repo: repo}
}

func (u *MeUser) Execute(in MeUserInput) (*MeUserOutput, error) {
	user, err := u.repo.GetByID(in.UserID)
	if err != nil {
		return nil, err
	}
	return &MeUserOutput{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Phone: user.Phone,
		CPF:   user.CPF,
	}, nil
}
