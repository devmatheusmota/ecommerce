package usecase

import (
	"github.com/ecommerce/services/users/internal/repository"
)

type DeleteAddress struct {
	addressRepository repository.AddressRepository
}

func NewDeleteAddress(addressRepository repository.AddressRepository) *DeleteAddress {
	return &DeleteAddress{addressRepository: addressRepository}
}

func (u *DeleteAddress) Execute(addressID, userID string) error {
	return u.addressRepository.Delete(addressID, userID)
}
