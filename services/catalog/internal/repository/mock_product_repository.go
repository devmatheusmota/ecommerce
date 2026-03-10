package repository

import "github.com/ecommerce/services/catalog/internal/domain"

type MockProductRepository struct {
	CreateFunc  func(*domain.Product) (*domain.Product, error)
	GetByIDFunc func(string) (*domain.Product, error)
	ListFunc    func(ListProductsFilter) ([]*domain.Product, int, error)
	UpdateFunc  func(*domain.Product) (*domain.Product, error)
	DeleteFunc  func(string) error
}

func (m *MockProductRepository) Create(product *domain.Product) (*domain.Product, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(product)
	}
	return product, nil
}

func (m *MockProductRepository) GetByID(id string) (*domain.Product, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, domain.ErrProductNotFound
}

func (m *MockProductRepository) List(filter ListProductsFilter) ([]*domain.Product, int, error) {
	if m.ListFunc != nil {
		return m.ListFunc(filter)
	}
	return nil, 0, nil
}

func (m *MockProductRepository) Update(product *domain.Product) (*domain.Product, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(product)
	}
	return product, nil
}

func (m *MockProductRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
