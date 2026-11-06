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
	// Часть имени — слова из букв, разделённые одиночными «-», «'» или пробелом
	// («Римский-Корсаков», «O'Neil», «Анна Мария»): начинается и заканчивается
	// буквой, двух разделителей подряд нет.
	prevSep := true // в начале строки разделитель запрещён
	for _, r := range v {
		switch {
		case unicode.IsLetter(r):
			prevSep = false
		case isNameSeparator(r):
			if prevSep {
				return "", ErrInvalidName
			}
			prevSep = true
		default:
			return "", ErrInvalidName
		}
	}
	if prevSep { // разделитель в конце строки
		return "", ErrInvalidName
	}
	return v, nil
}

func isNameSeparator(r rune) bool {
	return r == '-' || r == ' ' || r == '\''
}

func initial(s string) string {
	r, _ := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + "."
}
