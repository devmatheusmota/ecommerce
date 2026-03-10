package usecase

import (
	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/repository"
)

type ListAddresses struct {
	addressRepository repository.AddressRepository
}

func NewListAddresses(addressRepository repository.AddressRepository) *ListAddresses {
	return &ListAddresses{addressRepository: addressRepository}
}

func (u *ListAddresses) Execute(userID string) ([]*domain.Address, error) {
	return u.addressRepository.ListByUserID(userID)
}
