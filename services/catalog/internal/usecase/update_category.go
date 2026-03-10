package usecase

import (
	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/ecommerce/services/catalog/internal/repository"
)

type UpdateCategoryInput struct {
	ID       string
	Name     string
	Slug     string
	ParentID string
}

type UpdateCategory struct {
	categoryRepository repository.CategoryRepository
}

func NewUpdateCategory(categoryRepository repository.CategoryRepository) *UpdateCategory {
	return &UpdateCategory{categoryRepository: categoryRepository}
}

func (u *UpdateCategory) Execute(input UpdateCategoryInput) (*domain.Category, error) {
	category := &domain.Category{
		ID:       input.ID,
		Name:     input.Name,
		Slug:     input.Slug,
		ParentID: input.ParentID,
	}
	return u.categoryRepository.Update(category)
}
