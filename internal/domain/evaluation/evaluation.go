package evaluation

import (
	"time"

	"pos/internal/domain/rubric"
	"pos/internal/domain/shared"
)

type Evaluation struct {
	id          shared.ID
	studentID   shared.ID
	criterionID shared.ID
	score       Score
	comment     Comment
	gradedAt    time.Time
}

func Restore(id, studentID, criterionID shared.ID, score Score, comment Comment, gradedAt time.Time) *Evaluation {
	return &Evaluation{
		id:          id,
		studentID:   studentID,
		criterionID: criterionID,
		score:       score,
		comment:     comment,
		gradedAt:    gradedAt,
	}
}

func (e *Evaluation) ID() shared.ID          { return e.id }
func (e *Evaluation) StudentID() shared.ID   { return e.studentID }
func (e *Evaluation) CriterionID() shared.ID { return e.criterionID }
func (e *Evaluation) Score() Score           { return e.score }
func (e *Evaluation) Comment() Comment       { return e.comment }
func (e *Evaluation) GradedAt() time.Time    { return e.gradedAt }

func (e *Evaluation) Correct(c rubric.Criterion, score Score, comment Comment, at time.Time) error {
	if c.ID() != e.criterionID {
		return shared.ErrCriterionMismatch
	}
	if !c.MaxScore().Allows(score.Value()) {
		return shared.ErrScoreOutOfRange
	}
	e.score = score
	e.comment = comment
	e.gradedAt = at
	return nil
}
