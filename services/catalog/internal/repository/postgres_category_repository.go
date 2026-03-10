package repository

import (
	"database/sql"
	"errors"

	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/lib/pq"
)

type PostgresCategoryRepository struct {
	database *sql.DB
}

func NewPostgresCategoryRepository(database *sql.DB) *PostgresCategoryRepository {
	return &PostgresCategoryRepository{database: database}
}

func (r *PostgresCategoryRepository) Create(category *domain.Category) (*domain.Category, error) {
	var id, parentID sql.NullString
	if category.ParentID != "" {
		parentID = sql.NullString{String: category.ParentID, Valid: true}
	}

	var createdAt, updatedAt string
	err := r.database.QueryRow(`
		INSERT INTO categories (name, slug, parent_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`, category.Name, category.Slug, parentID).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, domain.ErrDuplicateSlug
		}
		return nil, err
	}
	category.ID = id.String
	category.CreatedAt = createdAt
	category.UpdatedAt = updatedAt
	return category, nil
}

func (r *PostgresCategoryRepository) GetByID(id string) (*domain.Category, error) {
	var category domain.Category
	var parentID sql.NullString
	err := r.database.QueryRow(`
		SELECT id, name, slug, parent_id, created_at, updated_at
		FROM categories WHERE id = $1
	`, id).Scan(&category.ID, &category.Name, &category.Slug, &parentID, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}
	if parentID.Valid {
		category.ParentID = parentID.String
	}
	return &category, nil
}

func (r *PostgresCategoryRepository) GetBySlug(slug string) (*domain.Category, error) {
	var category domain.Category
	var parentID sql.NullString
	err := r.database.QueryRow(`
		SELECT id, name, slug, parent_id, created_at, updated_at
		FROM categories WHERE slug = $1
	`, slug).Scan(&category.ID, &category.Name, &category.Slug, &parentID, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}
	if parentID.Valid {
		category.ParentID = parentID.String
	}
	return &category, nil
}

func (r *PostgresCategoryRepository) ListAll() ([]*domain.Category, error) {
	rows, err := r.database.Query(`
		SELECT id, name, slug, parent_id, created_at, updated_at
		FROM categories ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var category domain.Category
		var parentID sql.NullString
		if err := rows.Scan(&category.ID, &category.Name, &category.Slug, &parentID, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		if parentID.Valid {
			category.ParentID = parentID.String
		}
		categories = append(categories, &category)
	}
	return categories, rows.Err()
}

func (r *PostgresCategoryRepository) ListByParentID(parentID string) ([]*domain.Category, error) {
	query := `SELECT id, name, slug, parent_id, created_at, updated_at FROM categories WHERE `
	args := []any{}

	if parentID == "" {
		query += `parent_id IS NULL`
	} else {
		query += `parent_id = $1`
		args = append(args, parentID)
	}
	query += ` ORDER BY name`

	rows, err := r.database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var category domain.Category
		var scanParentID sql.NullString
		if err := rows.Scan(&category.ID, &category.Name, &category.Slug, &scanParentID, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		if scanParentID.Valid {
			category.ParentID = scanParentID.String
		}
		categories = append(categories, &category)
	}
	return categories, rows.Err()
}

func (r *PostgresCategoryRepository) Update(category *domain.Category) (*domain.Category, error) {
	var parentID sql.NullString
	if category.ParentID != "" {
		parentID = sql.NullString{String: category.ParentID, Valid: true}
	}

	result, err := r.database.Exec(`
		UPDATE categories SET name = $1, slug = $2, parent_id = $3, updated_at = now()
		WHERE id = $4
	`, category.Name, category.Slug, parentID, category.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, domain.ErrDuplicateSlug
		}
		return nil, err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, domain.ErrCategoryNotFound
	}
	return category, nil
}

func (r *PostgresCategoryRepository) Delete(id string) error {
	result, err := r.database.Exec(`DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrCategoryNotFound
	}
	return nil
}
