package repository

import (
	"sync"
	"time"
)

type MockPasswordResetRepository struct {
	mu     sync.RWMutex
	byHash map[string]*PasswordResetToken
}

func NewMockPasswordResetRepository() *MockPasswordResetRepository {
	return &MockPasswordResetRepository{byHash: make(map[string]*PasswordResetToken)}
}

func (m *MockPasswordResetRepository) Create(userID, tokenHash string, expiresAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.byHash[tokenHash] = &PasswordResetToken{ID: "token-1", UserID: userID, TokenHash: tokenHash, ExpiresAt: expiresAt}
	return nil
}

func (m *MockPasswordResetRepository) FindByTokenHash(tokenHash string) (*PasswordResetToken, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.byHash[tokenHash]
	if !ok {
		return nil, nil
	}
	clone := *t
	return &clone, nil
}

func (m *MockPasswordResetRepository) DeleteByID(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for h, t := range m.byHash {
		if t.ID == id {
			delete(m.byHash, h)
			break
		}
	}
	return nil
}
