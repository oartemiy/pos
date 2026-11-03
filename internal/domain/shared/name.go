package shared

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type Name struct {
	first      string
	last       string
	patronymic string
}

const maxNamePartLen = 64

func NewName(first, last, patronymic string) (Name, error) {
	l, err := normalizePart(last, true)
	if err != nil {
		return Name{}, err
	}
	f, err := normalizePart(first, true)
	if err != nil {
		return Name{}, err
	}
	p, err := normalizePart(patronymic, false)
	if err != nil {
		return Name{}, err
	}
	return Name{first: f, last: l, patronymic: p}, nil
}

func (n Name) First() string      { return n.first }
func (n Name) Last() string       { return n.last }
func (n Name) Patronymic() string { return n.patronymic }

func (n Name) String() string {
	parts := []string{n.last, n.first}
	if n.patronymic != "" {
		parts = append(parts, n.patronymic)
	}
	return strings.Join(parts, " ")
}

func (n Name) Short() string {
	s := n.last + " " + initial(n.first)
	if n.patronymic != "" {
		s += " " + initial(n.patronymic)
	}
	return s
}

func normalizePart(raw string, required bool) (string, error) {
	v := strings.Join(strings.Fields(raw), " ")
	if v == "" {
		if required {
			return "", ErrEmptyName
		}
		return "", nil
	}
	if utf8.RuneCountInString(v) > maxNamePartLen {
		return "", ErrNameTooLong
	}
	for _, r := range v {
		if !unicode.IsLetter(r) && r != '-' && r != ' ' && r != '\'' {
			return "", ErrInvalidName
		}
	}
	return v, nil
}

func initial(s string) string {
	r, _ := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + "."
}
