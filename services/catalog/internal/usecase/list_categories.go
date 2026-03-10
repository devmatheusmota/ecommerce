package usecase

import (
	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/ecommerce/services/catalog/internal/repository"
)

type ListCategories struct {
	categoryRepository repository.CategoryRepository
}

func NewListCategories(categoryRepository repository.CategoryRepository) *ListCategories {
	return &ListCategories{categoryRepository: categoryRepository}
}

func (u *ListCategories) Execute(parentID string) ([]*domain.Category, error) {
	if parentID != "" {
		return u.categoryRepository.ListByParentID(parentID)
	}
	return u.categoryRepository.ListAll()
}
