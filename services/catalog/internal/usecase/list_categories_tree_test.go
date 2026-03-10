package usecase

import (
	"testing"

	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/ecommerce/services/catalog/internal/repository"
)

func TestListCategoriesTree_Execute_Empty(t *testing.T) {
	mockRepository := &repository.MockCategoryRepository{
		ListAllFunc: func() ([]*domain.Category, error) {
			return []*domain.Category{}, nil
		},
	}
	listCategoriesTreeUseCase := NewListCategoriesTree(mockRepository)

	tree, err := listCategoriesTreeUseCase.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tree) != 0 {
		t.Errorf("expected empty tree, got %d roots", len(tree))
	}
}

func TestListCategoriesTree_Execute_Nested(t *testing.T) {
	mockRepository := &repository.MockCategoryRepository{
		ListAllFunc: func() ([]*domain.Category, error) {
			return []*domain.Category{
				{ID: "root-1", Name: "Root1", Slug: "root1", ParentID: ""},
				{ID: "child-1", Name: "Child1", Slug: "child1", ParentID: "root-1"},
				{ID: "child-2", Name: "Child2", Slug: "child2", ParentID: "root-1"},
				{ID: "grandchild", Name: "Grandchild", Slug: "grandchild", ParentID: "child-1"},
			}, nil
		},
	}
	listCategoriesTreeUseCase := NewListCategoriesTree(mockRepository)

	tree, err := listCategoriesTreeUseCase.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tree) != 1 {
		t.Fatalf("expected 1 root, got %d", len(tree))
	}
	root := tree[0]
	if root.Name != "Root1" {
		t.Errorf("expected root name Root1, got %s", root.Name)
	}
	if len(root.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(root.Children))
	}
	child1 := root.Children[0]
	if child1.Name != "Child1" {
		t.Errorf("expected Child1, got %s", child1.Name)
	}
	if len(child1.Children) != 1 {
		t.Fatalf("expected 1 grandchild, got %d", len(child1.Children))
	}
	if child1.Children[0].Name != "Grandchild" {
		t.Errorf("expected Grandchild, got %s", child1.Children[0].Name)
	}
}
