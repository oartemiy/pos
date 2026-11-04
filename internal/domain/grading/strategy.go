package grading

import "pos/internal/domain/rubric"

type Line struct {
	Criterion rubric.Criterion
	Score     float64
	Graded    bool
}

type Strategy interface {
	Name() string
	Calculate(lines []Line) float64
}