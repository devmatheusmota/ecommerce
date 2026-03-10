package usecase

import (
	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/repository"
)

type CreateAddressInput struct {
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

type CreateAddress struct {
	addressRepository repository.AddressRepository
}

func NewCreateAddress(addressRepository repository.AddressRepository) *CreateAddress {
	return &CreateAddress{addressRepository: addressRepository}
}

func (u *CreateAddress) Execute(in CreateAddressInput) (*domain.Address, error) {
	if in.IsDefaultBilling {
		_ = u.addressRepository.UnsetDefaultBillingForUser(in.UserID)
	}
	if in.IsDefaultShipping {
		_ = u.addressRepository.UnsetDefaultShippingForUser(in.UserID)
	}
	address := &domain.Address{
		UserID:            in.UserID,
		Street:            in.Street,
		Number:            in.Number,
		Complement:        in.Complement,
		Neighborhood:      in.Neighborhood,
		City:              in.City,
		State:             in.State,
		ZipCode:           in.ZipCode,
		Type:              in.Type,
		IsDefaultBilling:  in.IsDefaultBilling,
		IsDefaultShipping: in.IsDefaultShipping,
	}
	return u.addressRepository.Create(address)
}
