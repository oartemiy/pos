package rubric

import (
	"math"
	"pos/internal/domain/shared"
)

type MaxScore struct {
	value float64
}

func NewMaxScore(v float64) (MaxScore, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) || v <= 0 {
		return MaxScore{}, shared.ErrInvalidMaxScore
	}
	return MaxScore{value: v}, nil
}

func (m MaxScore) Value() float64 {
	return m.value
}

func (m MaxScore) Allows(score float64) bool {
	return score >= 0 && score <= m.value
}
