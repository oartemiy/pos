package dto

import "time"

type SetEvaluationInput struct {
	StudentID   string
	CriterionID string
	Score       float64
	Comment     string
}

type CorrectEvaluationInput struct {
	EvaluationID string
	Score        float64
	Comment      string
}

type EvaluationView struct {
	ID            string
	StudentID     string
	StudentName   string
	CriterionID   string
	CriterionName string
	Score         float64
	MaxScore      float64
	Comment       string
	GradedAt      time.Time
}
