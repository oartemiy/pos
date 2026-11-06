package shared

import (
	"net/mail"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	v := strings.ToLower(strings.TrimSpace(raw))
	if v == "" {
		return Email{}, ErrEmptyEmail
	}
	addr, err := mail.ParseAddress(v)
	if err != nil || addr.Address != v {
		return Email{}, ErrInvalidEmail
	}
	// ParseAddress следует RFC 5322 и пропускает адреса вроде «a@b»,
	// «ivan@localhost» и «a@[127.0.0.1]», поэтому домен проверяем отдельно.
	at := strings.LastIndexByte(v, '@')
	if !isValidDomain(v[at+1:]) {
		return Email{}, ErrInvalidEmail
	}
	return Email{v}, nil
}

func (e Email) String() string {
	return e.value
}

// isValidDomain требует хотя бы две метки через точку («edu.hse.ru»).
// Метка — буквы, цифры и дефисы, не с дефиса и не на дефис;
// домен верхнего уровня — только буквы, минимум две («ru», «рф»).
func isValidDomain(domain string) bool {
	labels := strings.Split(domain, ".")
	if len(labels) < 2 {
		return false
	}
	for _, l := range labels {
		if l == "" || strings.HasPrefix(l, "-") || strings.HasSuffix(l, "-") {
			return false
		}
		for _, r := range l {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
				return false
			}
		}
	}
	tld := labels[len(labels)-1]
	if utf8.RuneCountInString(tld) < 2 {
		return false
	}
	for _, r := range tld {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
