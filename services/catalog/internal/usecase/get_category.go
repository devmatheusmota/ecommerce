package usecase

import (
	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/ecommerce/services/catalog/internal/repository"
)

type GetCategory struct {
	categoryRepository repository.CategoryRepository
}

func NewGetCategory(categoryRepository repository.CategoryRepository) *GetCategory {
	return &GetCategory{categoryRepository: categoryRepository}
}

func (u *GetCategory) Execute(id string) (*domain.Category, error) {
	return u.categoryRepository.GetByID(id)
}
