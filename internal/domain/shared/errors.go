package shared

import "errors"

var (
	ErrEmptyEmail   error = errors.New("empty email")
	ErrInvalidEmail error = errors.New("invalid email")
	ErrEmptyName    error = errors.New("empty name")
	ErrNameTooLong  error = errors.New("name is too long")
	ErrInvalidName  error = errors.New("name consists of invalid characters")
	ErrIDParsing    error = errors.New("id parsing error")
)
