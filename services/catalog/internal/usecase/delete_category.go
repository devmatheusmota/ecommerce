package usecase

import (
	"github.com/ecommerce/services/catalog/internal/repository"
)

type DeleteCategory struct {
	categoryRepository repository.CategoryRepository
}

func NewDeleteCategory(categoryRepository repository.CategoryRepository) *DeleteCategory {
	return &DeleteCategory{categoryRepository: categoryRepository}
}

func (u *DeleteCategory) Execute(id string) error {
	return u.categoryRepository.Delete(id)
}
