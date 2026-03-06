package domain

import (
	"errors"
)

// Domain state / business rule errors (sentinel). Use errors.Is to check.
var (
	ErrDuplicateEmail     = errors.New("email already registered")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// ErrValidation is for invalid input (format, required fields). Used by the validation layer.
type ErrValidation string

func (e ErrValidation) Error() string { return string(e) }
