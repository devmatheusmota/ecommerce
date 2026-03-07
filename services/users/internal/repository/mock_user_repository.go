package repository

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/ecommerce/services/users/internal/domain"
)

// MockUserRepository is an in-memory UserRepository for testing.
// It is not safe for concurrent use by multiple goroutines that call Create;
// use from tests with a single goroutine or add locking as needed.
type MockUserRepository struct {
	mu      sync.RWMutex
	counter atomic.Uint64
	byID    map[string]*domain.User
	byEmail map[string]*domain.User
}

// NewMockUserRepository returns a new in-memory repository.
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		byID:    make(map[string]*domain.User),
		byEmail: make(map[string]*domain.User),
	}
}

// Create stores the user with a generated ID. Returns domain.ErrDuplicateEmail if email already exists.
func (m *MockUserRepository) Create(user *domain.User) (*domain.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.byEmail[user.Email]; exists {
		return nil, domain.ErrDuplicateEmail
	}
	user.ID = fmt.Sprintf("id-%d", m.counter.Add(1))
	clone := *user
	m.byID[user.ID] = &clone
	m.byEmail[user.Email] = &clone
	return &clone, nil
}

// GetByEmail returns the user by email or domain.ErrUserNotFound.
func (m *MockUserRepository) GetByEmail(email string) (*domain.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.byEmail[email]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	clone := *u
	return &clone, nil
}

// GetByID returns the user by ID or domain.ErrUserNotFound.
func (m *MockUserRepository) GetByID(id string) (*domain.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.byID[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	clone := *u
	return &clone, nil
}

// SetUser injects a user for testing (e.g. to simulate existing user for GetByID/GetByEmail).
func (m *MockUserRepository) SetUser(user *domain.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	clone := *user
	m.byID[user.ID] = &clone
	m.byEmail[user.Email] = &clone
}

// Reset clears all users. Useful between tests.
func (m *MockUserRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.byID = make(map[string]*domain.User)
	m.byEmail = make(map[string]*domain.User)
}

// ErrMock is used when the mock is configured to return an error (e.g. for GetByID failure).
var ErrMock = errors.New("mock error")
