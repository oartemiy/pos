package student

import "pos/internal/domain/shared"

type Student struct {
	id    shared.ID
	name  string
	email shared.Email
}

func NewStudent(email shared.Email) *Student {
	return &Student{shared.NewID(), "hell", email}
}

func (s *Student) Equal(other *Student) bool {
	return other != nil && s.id == other.id
}
