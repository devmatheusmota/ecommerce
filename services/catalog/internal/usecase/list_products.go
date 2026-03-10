package usecase

import (
	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/ecommerce/services/catalog/internal/repository"
)

type ListProductsFilter struct {
	SellerID   string
	CategoryID string
	Limit      int
	Offset     int
}

type ListProductsResult struct {
	Products []*domain.Product
	Total    int
}

type ListProducts struct {
	productRepository repository.ProductRepository
}

func NewListProducts(productRepository repository.ProductRepository) *ListProducts {
	return &ListProducts{productRepository: productRepository}
}

func (u *ListProducts) Execute(filter ListProductsFilter) (*ListProductsResult, error) {
	products, total, err := u.productRepository.List(repository.ListProductsFilter{
		SellerID:   filter.SellerID,
		CategoryID: filter.CategoryID,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
	})
	if err != nil {
		return nil, err
	}
	return &ListProductsResult{Products: products, Total: total}, nil
}
