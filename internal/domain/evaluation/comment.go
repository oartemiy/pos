package evaluation

import (
	"strings"
	"unicode/utf8"

	"pos/internal/domain/shared"
)

const maxCommentLen = 500

type Comment struct {
	value string
}

func NewComment(raw string) (Comment, error) {
	v := strings.TrimSpace(raw)
	if utf8.RuneCountInString(v) > maxCommentLen {
		return Comment{}, shared.ErrCommentTooLong
	}
	return Comment{value: v}, nil
}

func (c Comment) String() string { return c.value }
func (c Comment) IsEmpty() bool  { return c.value == "" }