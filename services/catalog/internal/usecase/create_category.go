package usecase

import (
	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/ecommerce/services/catalog/internal/repository"
)

type CreateCategoryInput struct {
	Name     string
	Slug     string
	ParentID string
}

type CreateCategory struct {
	categoryRepository repository.CategoryRepository
}

func NewCreateCategory(categoryRepository repository.CategoryRepository) *CreateCategory {
	return &CreateCategory{categoryRepository: categoryRepository}
}

func (u *CreateCategory) Execute(input CreateCategoryInput) (*domain.Category, error) {
	category := &domain.Category{
		Name:     input.Name,
		Slug:     input.Slug,
		ParentID: input.ParentID,
	}
	return u.categoryRepository.Create(category)
}
