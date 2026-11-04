package grading

import "pos/internal/domain/shared"

type Total struct {
	studentID shared.ID
	value     float64
	graded    int
	criteria  int
	strategy  string
}

func (t Total) StudentID() shared.ID { return t.studentID }
func (t Total) Value() float64       { return t.value }
func (t Total) Graded() int          { return t.graded }
func (t Total) Criteria() int        { return t.criteria }
func (t Total) Strategy() string     { return t.strategy }

func (t Total) IsComplete() bool { return t.graded == t.criteria }
