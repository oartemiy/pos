package student

import "pos/internal/domain/shared"

type Student struct {
	id    shared.ID
	name  shared.Name
	email shared.Email
}

func New(name shared.Name, email shared.Email) *Student {
	return &Student{shared.NewID(), name, email}
}

func (s *Student) Equal(other *Student) bool {
	return other != nil && s.id == other.id
}
func (s *Student) ID() shared.ID       { return s.id }
func (s *Student) Name() shared.Name   { return s.name }
func (s *Student) Email() shared.Email { return s.email }

func (s *Student) Rename(name shared.Name) {
	s.name = name
}

func (s *Student) ChangeEmail(email shared.Email) {
	s.email = email
}

func Restore(id shared.ID, name shared.Name, email shared.Email) *Student {
	return &Student{id: id, name: name, email: email}
}
