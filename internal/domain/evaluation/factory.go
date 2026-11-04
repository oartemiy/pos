package evaluation

import (
	"time"

	"pos/internal/domain/rubric"
	"pos/internal/domain/shared"
)

type Factory struct {
	now func() time.Time
}

func NewFactory(now func() time.Time) *Factory {
	return &Factory{now: now}
}

func (f *Factory) Create(
	rb *rubric.Rubric,
	studentID, criterionID shared.ID,
	rawScore float64,
	rawComment string,
) (*Evaluation, error) {
	c, err := rb.Criterion(criterionID)
	if err != nil {
		return nil, err
	}
	score, err := NewScore(rawScore)
	if err != nil {
		return nil, err
	}
	if !c.MaxScore().Allows(score.Value()) {
		return nil, shared.ErrScoreOutOfRange
	}
	comment, err := NewComment(rawComment)
	if err != nil {
		return nil, err
	}
	return &Evaluation{
		id:          shared.NewID(),
		studentID:   studentID,
		criterionID: c.ID(),
		score:       score,
		comment:     comment,
		gradedAt:    f.now(),
	}, nil
}
