package usecase

import (
	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/ecommerce/services/catalog/internal/repository"
)

type UpdateProductInput struct {
	ID          string
	SellerID    string
	CategoryID  string
	Title       string
	Description string
	Price       string
	Images      []string
}

type UpdateProduct struct {
	productRepository repository.ProductRepository
}

func NewUpdateProduct(productRepository repository.ProductRepository) *UpdateProduct {
	return &UpdateProduct{productRepository: productRepository}
}

func (u *UpdateProduct) Execute(input UpdateProductInput) (*domain.Product, error) {
	images := input.Images
	if images == nil {
		images = []string{}
	}
	product := &domain.Product{
		ID:          input.ID,
		SellerID:    input.SellerID,
		CategoryID:  input.CategoryID,
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		Images:      images,
	}
	return u.productRepository.Update(product)
}
