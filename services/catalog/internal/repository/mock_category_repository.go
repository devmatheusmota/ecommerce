package repository

import "github.com/ecommerce/services/catalog/internal/domain"

type MockCategoryRepository struct {
	CreateFunc         func(*domain.Category) (*domain.Category, error)
	GetByIDFunc        func(string) (*domain.Category, error)
	GetBySlugFunc      func(string) (*domain.Category, error)
	ListAllFunc        func() ([]*domain.Category, error)
	ListByParentIDFunc func(string) ([]*domain.Category, error)
	UpdateFunc         func(*domain.Category) (*domain.Category, error)
	DeleteFunc         func(string) error
}

func (m *MockCategoryRepository) Create(category *domain.Category) (*domain.Category, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(category)
	}
	return category, nil
}

func (m *MockCategoryRepository) GetByID(id string) (*domain.Category, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, domain.ErrCategoryNotFound
}

func (m *MockCategoryRepository) GetBySlug(slug string) (*domain.Category, error) {
	if m.GetBySlugFunc != nil {
		return m.GetBySlugFunc(slug)
	}
	return nil, domain.ErrCategoryNotFound
}

func (m *MockCategoryRepository) ListAll() ([]*domain.Category, error) {
	if m.ListAllFunc != nil {
		return m.ListAllFunc()
	}
	return nil, nil
}

func (m *MockCategoryRepository) ListByParentID(parentID string) ([]*domain.Category, error) {
	if m.ListByParentIDFunc != nil {
		return m.ListByParentIDFunc(parentID)
	}
	return nil, nil
}

func (m *MockCategoryRepository) Update(category *domain.Category) (*domain.Category, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(category)
	}
	return category, nil
}

func (m *MockCategoryRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
