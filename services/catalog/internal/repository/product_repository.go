package repository

import "github.com/ecommerce/services/catalog/internal/domain"

type ListProductsFilter struct {
	SellerID        string
	CategoryID      string
	CategoryIDs     []string // exact match: category_id IN (...); used for related products
	ExcludeProductID string
	Limit           int
	Offset          int
}

type ProductRepository interface {
	Create(product *domain.Product) (*domain.Product, error)
	GetByID(id string) (*domain.Product, error)
	List(filter ListProductsFilter) ([]*domain.Product, int, error)
	Update(product *domain.Product) (*domain.Product, error)
	Delete(id string) error
}
