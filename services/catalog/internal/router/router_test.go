package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/ecommerce/services/catalog/internal/repository"
)

func TestHealth(t *testing.T) {
	handler := NewWithRepositories(
		&repository.MockCategoryRepository{},
		&repository.MockProductRepository{},
	)

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", response.Code)
	}
}

func TestListCategories_Empty(t *testing.T) {
	mockRepository := &repository.MockCategoryRepository{
		ListAllFunc: func() ([]*domain.Category, error) {
			return []*domain.Category{}, nil
		},
	}
	handler := NewWithRepositories(mockRepository, &repository.MockProductRepository{})

	request := httptest.NewRequest(http.MethodGet, "/v1/categories", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", response.Code)
	}
}

func TestListCategoriesTree(t *testing.T) {
	mockRepository := &repository.MockCategoryRepository{
		ListAllFunc: func() ([]*domain.Category, error) {
			return []*domain.Category{
				{ID: "r1", Name: "Root", Slug: "root", ParentID: ""},
			}, nil
		},
	}
	handler := NewWithRepositories(mockRepository, &repository.MockProductRepository{})

	request := httptest.NewRequest(http.MethodGet, "/v1/categories/tree", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", response.Code)
	}
}
