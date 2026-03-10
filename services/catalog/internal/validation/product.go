package validation

import (
	"regexp"
	"strings"

	"github.com/ecommerce/services/catalog/internal/domain"
)

var priceRegex = regexp.MustCompile(`^\d+(\.\d{1,2})?$`)

type ProductInput struct {
	SellerID    string
	CategoryID  string
	Title       string
	Description string
	Price       string
	Images      []string
}

func ValidateProductInput(input *ProductInput) error {
	sellerID := strings.TrimSpace(input.SellerID)
	if sellerID == "" {
		return domain.ErrValidation("seller_id is required")
	}
	input.SellerID = sellerID

	categoryID := strings.TrimSpace(input.CategoryID)
	if categoryID == "" {
		return domain.ErrValidation("category_id is required")
	}
	input.CategoryID = categoryID

	title := strings.TrimSpace(input.Title)
	if title == "" {
		return domain.ErrValidation("title is required")
	}
	if len(title) > 500 {
		return domain.ErrValidation("title must be at most 500 characters")
	}
	input.Title = title

	description := strings.TrimSpace(input.Description)
	if description == "" {
		return domain.ErrValidation("description is required")
	}
	input.Description = description

	price := strings.TrimSpace(input.Price)
	if price == "" {
		return domain.ErrValidation("price is required")
	}
	if !priceRegex.MatchString(price) {
		return domain.ErrValidation("price must be a positive number with up to 2 decimal places (e.g. 99.90)")
	}
	input.Price = price

	for i, image := range input.Images {
		input.Images[i] = strings.TrimSpace(image)
	}
	return nil
}
