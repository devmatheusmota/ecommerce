package usecase

import (
	"errors"
	"testing"

	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/ecommerce/services/catalog/internal/repository"
)

func TestCreateCategory_Execute_Success(t *testing.T) {
	mockRepository := &repository.MockCategoryRepository{
		CreateFunc: func(category *domain.Category) (*domain.Category, error) {
			category.ID = "cat-123"
			return category, nil
		},
	}
	createCategoryUseCase := NewCreateCategory(mockRepository)

	category, err := createCategoryUseCase.Execute(CreateCategoryInput{
		Name:     "Eletrônicos",
		Slug:     "eletronicos",
		ParentID: "",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if category.ID != "cat-123" {
		t.Errorf("expected ID cat-123, got %s", category.ID)
	}
	if category.Name != "Eletrônicos" {
		t.Errorf("expected Name Eletrônicos, got %s", category.Name)
	}
	if category.Slug != "eletronicos" {
		t.Errorf("expected Slug eletronicos, got %s", category.Slug)
	}
}

func TestCreateCategory_Execute_DuplicateSlug(t *testing.T) {
	mockRepository := &repository.MockCategoryRepository{
		CreateFunc: func(*domain.Category) (*domain.Category, error) {
			return nil, domain.ErrDuplicateSlug
		},
	}
	createCategoryUseCase := NewCreateCategory(mockRepository)

	_, err := createCategoryUseCase.Execute(CreateCategoryInput{
		Name: "Eletrônicos",
		Slug: "eletronicos",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrDuplicateSlug) {
		t.Errorf("expected ErrDuplicateSlug, got %v", err)
	}
}
