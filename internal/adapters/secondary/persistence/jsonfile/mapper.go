package jsonfile

import (
	"fmt"
	"pos/internal/domain/shared"
	"pos/internal/domain/student"
)

type studentRecord struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Patronymic string `json:"patronymic,omitempty"`
	Email      string `json:"email"`
}

func toStudentRecord(s *student.Student) studentRecord {
	n := s.Name()
	return studentRecord{
		ID:         s.ID().String(),
		LastName:   n.Last(),
		FirstName:  n.First(),
		Patronymic: n.Patronymic(),
		Email:      s.Email().String(),
	}
}

func toStudent(r studentRecord) (*student.Student, error) {
	id, err := shared.ParseID(r.ID)
	if err != nil {
		return nil, fmt.Errorf("jsonfile: student %q: %w", r.ID, err)
	}
	name, err := shared.NewName(r.LastName, r.FirstName, r.Patronymic)
	if err != nil {
		return nil, fmt.Errorf("jsonfile: student %q: %w", r.ID, err)
	}
	email, err := shared.NewEmail(r.Email)
	if err != nil {
		return nil, fmt.Errorf("jsonfile: student %q: %w", r.ID, err)
	}
	return student.Restore(id, name, email), nil
}
