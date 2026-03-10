package usecase

import (
	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/ecommerce/services/catalog/internal/repository"
)

type GetProduct struct {
	productRepository repository.ProductRepository
}

func NewGetProduct(productRepository repository.ProductRepository) *GetProduct {
	return &GetProduct{productRepository: productRepository}
}

func (u *GetProduct) Execute(id string) (*domain.Product, error) {
	return u.productRepository.GetByID(id)
}
