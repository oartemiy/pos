package shared

import "errors"

var (
	ErrEmptyEmail   error = errors.New("empty email")
	ErrInvalidEmail error = errors.New("invalid email")
)
