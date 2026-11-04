package grading

import (
	"pos/internal/domain/evaluation"
	"pos/internal/domain/rubric"
	"pos/internal/domain/shared"
)

// Service — Domain Service расчёта итога.
type Service struct {
	strategy Strategy
}

func NewService(s Strategy) *Service {
	return &Service{strategy: s}
}

func (s *Service) Calculate(studentID shared.ID, rb *rubric.Rubric, evals []*evaluation.Evaluation) (Total, error) {
	criteria := rb.Criteria()
	if len(criteria) == 0 {
		return Total{}, shared.ErrEmptyRubric
	}

	scores := make(map[shared.ID]float64, len(evals))
	for _, e := range evals {
		if e.StudentID() != studentID {
			return Total{}, shared.ErrForeignEvaluation
		}
		if _, dup := scores[e.CriterionID()]; dup {
			return Total{}, shared.ErrDuplicateEvaluation
		}
		scores[e.CriterionID()] = e.Score().Value()
	}

	lines := make([]Line, 0, len(criteria))
	graded := 0
	for _, c := range criteria {
		score, ok := scores[c.ID()]
		if ok {
			graded++
		}
		lines = append(lines, Line{Criterion: c, Score: score, Graded: ok})
	}

	return Total{
		studentID: studentID,
		value:     s.strategy.Calculate(lines),
		graded:    graded,
		criteria:  len(criteria),
		strategy:  s.strategy.Name(),
	}, nil
}