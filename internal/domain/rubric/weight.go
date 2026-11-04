package rubric

import (
	"math"

	"pos/internal/domain/shared"
)

type Weight struct {
	value float64
}

func NewWeight(v float64) (Weight, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) || v <= 0 {
		return Weight{}, shared.ErrInvalidWeight
	}
	return Weight{value: v}, nil
}

func (w Weight) Value() float64 { return w.value }