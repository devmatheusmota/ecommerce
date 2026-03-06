package repository

import "github.com/ecommerce/services/users/internal/domain"

type UserRepository interface {
	Create(user *domain.User) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
}
