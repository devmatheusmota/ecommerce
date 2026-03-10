package repository

import (
	"database/sql"
	"errors"

	"github.com/ecommerce/services/users/internal/domain"
)

type PostgresAddressRepository struct {
	db *sql.DB
}

func NewPostgresAddressRepository(db *sql.DB) *PostgresAddressRepository {
	return &PostgresAddressRepository{db: db}
}

func (r *PostgresAddressRepository) Create(address *domain.Address) (*domain.Address, error) {
	var id string
	err := r.db.QueryRow(`
		INSERT INTO addresses (user_id, street, number, complement, neighborhood, city, state, zip_code, type, is_default_billing, is_default_shipping)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`,
		address.UserID, address.Street, address.Number, address.Complement, address.Neighborhood,
		address.City, address.State, address.ZipCode, address.Type, address.IsDefaultBilling, address.IsDefaultShipping,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	address.ID = id
	return address, nil
}

func (r *PostgresAddressRepository) GetByIDAndUserID(id, userID string) (*domain.Address, error) {
	var a domain.Address
	err := r.db.QueryRow(`
		SELECT id, user_id, street, number, complement, neighborhood, city, state, zip_code, type, is_default_billing, is_default_shipping
		FROM addresses WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(&a.ID, &a.UserID, &a.Street, &a.Number, &a.Complement, &a.Neighborhood, &a.City, &a.State, &a.ZipCode, &a.Type, &a.IsDefaultBilling, &a.IsDefaultShipping)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAddressNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *PostgresAddressRepository) ListByUserID(userID string) ([]*domain.Address, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, street, number, complement, neighborhood, city, state, zip_code, type, is_default_billing, is_default_shipping
		FROM addresses WHERE user_id = $1 ORDER BY created_at
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.Address
	for rows.Next() {
		var a domain.Address
		if err := rows.Scan(&a.ID, &a.UserID, &a.Street, &a.Number, &a.Complement, &a.Neighborhood, &a.City, &a.State, &a.ZipCode, &a.Type, &a.IsDefaultBilling, &a.IsDefaultShipping); err != nil {
			return nil, err
		}
		list = append(list, &a)
	}
	return list, rows.Err()
}

func (r *PostgresAddressRepository) Update(address *domain.Address) (*domain.Address, error) {
	result, err := r.db.Exec(`
		UPDATE addresses SET street = $1, number = $2, complement = $3, neighborhood = $4, city = $5, state = $6, zip_code = $7, type = $8, is_default_billing = $9, is_default_shipping = $10
		WHERE id = $11 AND user_id = $12
	`, address.Street, address.Number, address.Complement, address.Neighborhood, address.City, address.State, address.ZipCode, address.Type, address.IsDefaultBilling, address.IsDefaultShipping, address.ID, address.UserID)
	if err != nil {
		return nil, err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, domain.ErrAddressNotFound
	}
	return address, nil
}

func (r *PostgresAddressRepository) Delete(id, userID string) error {
	result, err := r.db.Exec(`DELETE FROM addresses WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrAddressNotFound
	}
	return nil
}

func (r *PostgresAddressRepository) UnsetDefaultBillingForUser(userID string) error {
	_, err := r.db.Exec(`UPDATE addresses SET is_default_billing = false WHERE user_id = $1`, userID)
	return err
}

func (r *PostgresAddressRepository) UnsetDefaultShippingForUser(userID string) error {
	_, err := r.db.Exec(`UPDATE addresses SET is_default_shipping = false WHERE user_id = $1`, userID)
	return err
}
