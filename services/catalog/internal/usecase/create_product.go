package usecase

import (
	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/ecommerce/services/catalog/internal/repository"
)

type CreateProductInput struct {
	SellerID    string
	CategoryID  string
	Title       string
	Description string
	Price       string
	Images      []string
}

type CreateProduct struct {
	productRepository repository.ProductRepository
}

func NewCreateProduct(productRepository repository.ProductRepository) *CreateProduct {
	return &CreateProduct{productRepository: productRepository}
}

func (u *CreateProduct) Execute(input CreateProductInput) (*domain.Product, error) {
	images := input.Images
	if images == nil {
		images = []string{}
	}
	product := &domain.Product{
		SellerID:    input.SellerID,
		CategoryID:  input.CategoryID,
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		Images:      images,
	}
	return u.productRepository.Create(product)
}
