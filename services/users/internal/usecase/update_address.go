package usecase

import (
	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/repository"
)

type UpdateAddressInput struct {
	AddressID         string
	UserID            string
	Street            string
	Number            string
	Complement        string
	Neighborhood      string
	City              string
	State             string
	ZipCode           string
	Type              string
	IsDefaultBilling  bool
	IsDefaultShipping bool
}

type UpdateAddress struct {
	addressRepository repository.AddressRepository
}

func NewUpdateAddress(addressRepository repository.AddressRepository) *UpdateAddress {
	return &UpdateAddress{addressRepository: addressRepository}
}

func (u *UpdateAddress) Execute(in UpdateAddressInput) (*domain.Address, error) {
	address, err := u.addressRepository.GetByIDAndUserID(in.AddressID, in.UserID)
	if err != nil {
		return nil, err
	}
	if in.IsDefaultBilling {
		_ = u.addressRepository.UnsetDefaultBillingForUser(in.UserID)
	}
	if in.IsDefaultShipping {
		_ = u.addressRepository.UnsetDefaultShippingForUser(in.UserID)
	}
	address.Street = in.Street
	address.Number = in.Number
	address.Complement = in.Complement
	address.Neighborhood = in.Neighborhood
	address.City = in.City
	address.State = in.State
	address.ZipCode = in.ZipCode
	address.Type = in.Type
	address.IsDefaultBilling = in.IsDefaultBilling
	address.IsDefaultShipping = in.IsDefaultShipping
	return u.addressRepository.Update(address)
}
