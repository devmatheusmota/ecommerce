package validation

import (
	"regexp"
	"strings"

	"github.com/ecommerce/services/catalog/internal/domain"
)

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type CategoryInput struct {
	Name     string
	Slug     string
	ParentID string
}

func ValidateCategoryInput(input *CategoryInput) error {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return domain.ErrValidation("name is required")
	}
	if len(name) > 255 {
		return domain.ErrValidation("name must be at most 255 characters")
	}
	input.Name = name

	slug := strings.TrimSpace(strings.ToLower(input.Slug))
	if slug == "" {
		return domain.ErrValidation("slug is required")
	}
	if len(slug) > 255 {
		return domain.ErrValidation("slug must be at most 255 characters")
	}
	if !slugRegex.MatchString(slug) {
		return domain.ErrValidation("slug must contain only lowercase letters, numbers and hyphens (e.g. eletronicos)")
	}
	input.Slug = slug
	input.ParentID = strings.TrimSpace(input.ParentID)
	return nil
}
