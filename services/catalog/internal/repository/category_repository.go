package repository

import "github.com/ecommerce/services/catalog/internal/domain"

type CategoryRepository interface {
	Create(category *domain.Category) (*domain.Category, error)
	GetByID(id string) (*domain.Category, error)
	GetBySlug(slug string) (*domain.Category, error)
	ListAll() ([]*domain.Category, error)
	ListByParentID(parentID string) ([]*domain.Category, error)
	Update(category *domain.Category) (*domain.Category, error)
	Delete(id string) error
}
