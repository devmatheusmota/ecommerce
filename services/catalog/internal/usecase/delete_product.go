package usecase

import (
	"github.com/ecommerce/services/catalog/internal/repository"
)

type DeleteProduct struct {
	productRepository repository.ProductRepository
}

func NewDeleteProduct(productRepository repository.ProductRepository) *DeleteProduct {
	return &DeleteProduct{productRepository: productRepository}
}

func (u *DeleteProduct) Execute(id string) error {
	return u.productRepository.Delete(id)
}
