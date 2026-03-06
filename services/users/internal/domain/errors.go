package domain

import (
	"errors"
)

var (
	ErrDuplicateEmail = errors.New("email already registered")
)

type ErrValidation string

func (e ErrValidation) Error() string { return string(e) }
