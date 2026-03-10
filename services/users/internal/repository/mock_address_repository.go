package repository

import (
	"fmt"
	"sync"

	"github.com/ecommerce/services/users/internal/domain"
)

type MockAddressRepository struct {
	mu     sync.RWMutex
	byID   map[string]*domain.Address
	byUser map[string][]*domain.Address
}

func NewMockAddressRepository() *MockAddressRepository {
	return &MockAddressRepository{
		byID:   make(map[string]*domain.Address),
		byUser: make(map[string][]*domain.Address),
	}
}

func (m *MockAddressRepository) Create(address *domain.Address) (*domain.Address, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	address.ID = fmt.Sprintf("addr-%d", len(m.byID)+1)
	clone := *address
	m.byID[address.ID] = &clone
	m.byUser[address.UserID] = append(m.byUser[address.UserID], &clone)
	return &clone, nil
}

func (m *MockAddressRepository) GetByIDAndUserID(id, userID string) (*domain.Address, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	a, ok := m.byID[id]
	if !ok || a.UserID != userID {
		return nil, domain.ErrAddressNotFound
	}
	clone := *a
	return &clone, nil
}

func (m *MockAddressRepository) ListByUserID(userID string) ([]*domain.Address, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	list := m.byUser[userID]
	out := make([]*domain.Address, len(list))
	for i, a := range list {
		clone := *a
		out[i] = &clone
	}
	return out, nil
}

func (m *MockAddressRepository) Update(address *domain.Address) (*domain.Address, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.byID[address.ID]
	if !ok || a.UserID != address.UserID {
		return nil, domain.ErrAddressNotFound
	}
	*a = *address
	clone := *a
	return &clone, nil
}

func (m *MockAddressRepository) Delete(id, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.byID[id]
	if !ok || a.UserID != userID {
		return domain.ErrAddressNotFound
	}
	delete(m.byID, id)
	list := m.byUser[userID]
	for i, ad := range list {
		if ad.ID == id {
			m.byUser[userID] = append(list[:i], list[i+1:]...)
			break
		}
	}
	return nil
}

func (m *MockAddressRepository) UnsetDefaultBillingForUser(userID string) error {
	for _, a := range m.byID {
		if a.UserID == userID {
			a.IsDefaultBilling = false
		}
	}
	return nil
}

func (m *MockAddressRepository) UnsetDefaultShippingForUser(userID string) error {
	for _, a := range m.byID {
		if a.UserID == userID {
			a.IsDefaultShipping = false
		}
	}
	return nil
}
