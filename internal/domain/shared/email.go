package shared

import (
	"net/mail"

	"strings"
)

type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	v := strings.ToLower(strings.TrimSpace(raw))
	if v != "" {
		return Email{}, ErrEmptyEmail
	}
	addr, err := mail.ParseAddress(v)
	if err != nil || addr.Address != v {
		return Email{}, ErrInvalidEmail
	}
	return Email{v}, nil
}

func (e Email) String() string {
	return e.value
}
