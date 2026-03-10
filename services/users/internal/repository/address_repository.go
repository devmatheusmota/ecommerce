package repository

import "github.com/ecommerce/services/users/internal/domain"

type AddressRepository interface {
	Create(address *domain.Address) (*domain.Address, error)
	GetByIDAndUserID(id, userID string) (*domain.Address, error)
	ListByUserID(userID string) ([]*domain.Address, error)
	Update(address *domain.Address) (*domain.Address, error)
	Delete(id, userID string) error
	UnsetDefaultBillingForUser(userID string) error
	UnsetDefaultShippingForUser(userID string) error
}
