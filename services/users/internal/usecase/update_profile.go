package usecase

import (
	"strings"

	"github.com/ecommerce/services/users/internal/repository"
)

type UpdateProfileInput struct {
	UserID string
	Name   string
	Phone  string
	CPF    string
}

type UpdateProfile struct {
	repository repository.UserRepository
}

func NewUpdateProfile(repository repository.UserRepository) *UpdateProfile {
	return &UpdateProfile{repository: repository}
}

func (u *UpdateProfile) Execute(in UpdateProfileInput) (*MeUserOutput, error) {
	user, err := u.repository.GetByID(in.UserID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(in.Name) != "" {
		user.Name = strings.TrimSpace(in.Name)
	}
	if strings.TrimSpace(in.Phone) != "" {
		user.Phone = strings.TrimSpace(in.Phone)
	}
	if strings.TrimSpace(in.CPF) != "" {
		user.CPF = strings.TrimSpace(in.CPF)
	}

	updated, err := u.repository.Update(user)
	if err != nil {
		return nil, err
	}

	return &MeUserOutput{
		ID:    updated.ID,
		Email: updated.Email,
		Name:  updated.Name,
		Phone: updated.Phone,
		CPF:   updated.CPF,
	}, nil
}
