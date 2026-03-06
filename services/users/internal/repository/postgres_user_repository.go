package repository

import (
	"database/sql"
	"errors"

	"github.com/ecommerce/services/users/internal/domain"
	"github.com/lib/pq"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(user *domain.User) (*domain.User, error) {
	var id string
	err := r.db.QueryRow(`
		INSERT INTO users (email, name, phone, cpf, password_hash)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, user.Email, user.Name, user.Phone, user.CPF, user.PasswordHash).Scan(&id)

	if err != nil {
		if isUniqueViolation(err) {
			return nil, domain.ErrDuplicateEmail
		}
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (r *PostgresUserRepository) GetByEmail(email string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(`
		SELECT id, email, name, phone, cpf, password_hash
		FROM users WHERE email = $1
	`, email).Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.CPF, &u.PasswordHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}
