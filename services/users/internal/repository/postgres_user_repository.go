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

func (r *PostgresUserRepository) GetByID(id string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(`
		SELECT id, email, name, phone, cpf, password_hash
		FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.Email, &u.Name, &u.Phone, &u.CPF, &u.PasswordHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) Update(user *domain.User) (*domain.User, error) {
	result, err := r.db.Exec(`
		UPDATE users SET name = $1, phone = $2, cpf = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`, user.Name, user.Phone, user.CPF, user.ID)
	if err != nil {
		return nil, err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *PostgresUserRepository) UpdatePassword(userID, passwordHash string) error {
	result, err := r.db.Exec(`UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, passwordHash, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
