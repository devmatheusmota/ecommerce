package usecase

import (
	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/repository"
)

type GetAddress struct {
	addressRepository repository.AddressRepository
}

func NewGetAddress(addressRepository repository.AddressRepository) *GetAddress {
	return &GetAddress{addressRepository: addressRepository}
}

func (u *GetAddress) Execute(addressID, userID string) (*domain.Address, error) {
	return u.addressRepository.GetByIDAndUserID(addressID, userID)
}
