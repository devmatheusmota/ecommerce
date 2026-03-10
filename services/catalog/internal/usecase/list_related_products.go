package usecase

import (
	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/ecommerce/services/catalog/internal/repository"
)

type ListRelatedProducts struct {
	productRepository  repository.ProductRepository
	categoryRepository repository.CategoryRepository
}

func NewListRelatedProducts(productRepository repository.ProductRepository, categoryRepository repository.CategoryRepository) *ListRelatedProducts {
	return &ListRelatedProducts{
		productRepository:  productRepository,
		categoryRepository: categoryRepository,
	}
}

// Execute returns products from the same category or the parent category, excluding the given product. Limit 12.
func (u *ListRelatedProducts) Execute(productID string) ([]*domain.Product, error) {
	product, err := u.productRepository.GetByID(productID)
	if err != nil {
		return nil, err
	}
	category, err := u.categoryRepository.GetByID(product.CategoryID)
	if err != nil {
		return nil, err
	}
	categoryIDs := []string{product.CategoryID}
	if category.ParentID != "" {
		categoryIDs = append(categoryIDs, category.ParentID)
	}
	products, _, err := u.productRepository.List(repository.ListProductsFilter{
		CategoryIDs:      categoryIDs,
		ExcludeProductID: productID,
		Limit:            12,
		Offset:           0,
	})
	if err != nil {
		return nil, err
	}
	return products, nil
}
