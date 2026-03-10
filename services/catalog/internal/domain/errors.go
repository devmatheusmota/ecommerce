package domain

import "errors"

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrProductNotFound  = errors.New("product not found")
	ErrDuplicateSlug    = errors.New("category slug already exists")
)

type ErrValidation string

func (e ErrValidation) Error() string { return string(e) }
