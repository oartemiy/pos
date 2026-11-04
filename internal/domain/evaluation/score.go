package evaluation

import (
	"math"

	"pos/internal/domain/shared"
)

type Score struct {
	value float64
}

func NewScore(v float64) (Score, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) || v < 0 {
		return Score{}, shared.ErrInvalidScore
	}
	return Score{value: v}, nil
}

func (s Score) Value() float64 { return s.value }
