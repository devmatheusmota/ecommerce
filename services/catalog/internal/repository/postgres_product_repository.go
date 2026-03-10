package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ecommerce/services/catalog/internal/domain"
)

type PostgresProductRepository struct {
	database *sql.DB
}

func NewPostgresProductRepository(database *sql.DB) *PostgresProductRepository {
	return &PostgresProductRepository{database: database}
}

func (r *PostgresProductRepository) Create(product *domain.Product) (*domain.Product, error) {
	imagesJSON, err := json.Marshal(product.Images)
	if err != nil {
		return nil, err
	}
	if len(product.Images) == 0 {
		imagesJSON = []byte("[]")
	}

	var id, createdAt, updatedAt string
	err = r.database.QueryRow(`
		INSERT INTO products (seller_id, category_id, title, description, price, images)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`, product.SellerID, product.CategoryID, product.Title, product.Description, product.Price, imagesJSON).
		Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	product.ID = id
	product.CreatedAt = createdAt
	product.UpdatedAt = updatedAt
	return product, nil
}

func (r *PostgresProductRepository) GetByID(id string) (*domain.Product, error) {
	var product domain.Product
	var imagesJSON []byte
	err := r.database.QueryRow(`
		SELECT id, seller_id, category_id, title, description, price, images, created_at, updated_at
		FROM products WHERE id = $1
	`, id).Scan(&product.ID, &product.SellerID, &product.CategoryID, &product.Title, &product.Description,
		&product.Price, &imagesJSON, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}
	if len(imagesJSON) > 0 {
		_ = json.Unmarshal(imagesJSON, &product.Images)
	}
	return &product, nil
}

func (r *PostgresProductRepository) List(filter ListProductsFilter) ([]*domain.Product, int, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	where := "1=1"
	args := []any{}
	placeholder := 1

	if filter.SellerID != "" {
		where += fmt.Sprintf(" AND seller_id = $%d", placeholder)
		args = append(args, filter.SellerID)
		placeholder++
	}
	if filter.ExcludeProductID != "" {
		where += fmt.Sprintf(" AND id != $%d", placeholder)
		args = append(args, filter.ExcludeProductID)
		placeholder++
	}
	useCategoryIDs := len(filter.CategoryIDs) > 0
	if useCategoryIDs {
		placeholders := make([]string, len(filter.CategoryIDs))
		for i, categoryID := range filter.CategoryIDs {
			placeholders[i] = fmt.Sprintf("$%d", placeholder)
			args = append(args, categoryID)
			placeholder++
		}
		where += " AND category_id IN (" + strings.Join(placeholders, ",") + ")"
	}
	categoryPlaceholder := placeholder
	if filter.CategoryID != "" && !useCategoryIDs {
		args = append(args, filter.CategoryID)
		categoryPlaceholder = placeholder
		placeholder++
	}

	// When filtering by single category, include category and all descendants (e.g. "Casa e Construção" includes "Decoração", "Móveis").
	withCategoryCTE := filter.CategoryID != "" && !useCategoryIDs
	var categoryCondition string
	if withCategoryCTE {
		categoryCondition = " AND category_id IN (SELECT id FROM category_and_descendants)"
	} else {
		categoryCondition = ""
	}
	recursiveCTE := ""
	if withCategoryCTE {
		recursiveCTE = fmt.Sprintf(
			`WITH RECURSIVE category_and_descendants AS (
				SELECT id FROM categories WHERE id = $%d
				UNION ALL
				SELECT c.id FROM categories c INNER JOIN category_and_descendants d ON c.parent_id = d.id
			) `,
			categoryPlaceholder,
		)
	}

	countQuery := recursiveCTE + `SELECT COUNT(*) FROM products WHERE ` + where + categoryCondition
	var total int
	err := r.database.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	countArgs := len(args)
	args = append(args, limit, offset)
	listQuery := recursiveCTE + `SELECT id, seller_id, category_id, title, description, price, images, created_at, updated_at
		FROM products WHERE ` + where + categoryCondition + fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, countArgs+1, countArgs+2)
	rows, err := r.database.Query(listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		var imagesJSON []byte
		if err := rows.Scan(&product.ID, &product.SellerID, &product.CategoryID, &product.Title, &product.Description,
			&product.Price, &imagesJSON, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if len(imagesJSON) > 0 {
			_ = json.Unmarshal(imagesJSON, &product.Images)
		}
		products = append(products, &product)
	}
	return products, total, rows.Err()
}

func (r *PostgresProductRepository) Update(product *domain.Product) (*domain.Product, error) {
	imagesJSON, err := json.Marshal(product.Images)
	if err != nil {
		return nil, err
	}
	if len(product.Images) == 0 {
		imagesJSON = []byte("[]")
	}

	result, err := r.database.Exec(`
		UPDATE products SET seller_id = $1, category_id = $2, title = $3, description = $4, price = $5, images = $6, updated_at = now()
		WHERE id = $7
	`, product.SellerID, product.CategoryID, product.Title, product.Description, product.Price, imagesJSON, product.ID)
	if err != nil {
		return nil, err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, domain.ErrProductNotFound
	}
	return product, nil
}

func (r *PostgresProductRepository) Delete(id string) error {
	result, err := r.database.Exec(`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}
