package facade

import (
	"context"
	"pos/internal/application/dto"
	"pos/internal/application/usecase"
)

type Facade struct {
	students    *usecase.Students
	criteria    *usecase.Criteria
	evaluations *usecase.Evaluations
	totals      *usecase.Totals
}

func New(
	s *usecase.Students,
	c *usecase.Criteria,
	e *usecase.Evaluations,
	t *usecase.Totals,
) *Facade {
	return &Facade{students: s, criteria: c, evaluations: e, totals: t}
}

// Студенты
func (f *Facade) AddStudent(ctx context.Context, in dto.AddStudentInput) (dto.StudentView, error) {
	return f.students.Add(ctx, in)
}

func (f *Facade) ListStudents(ctx context.Context) ([]dto.StudentView, error) {
	return f.students.List(ctx)
}

func (f *Facade) RemoveStudent(ctx context.Context, id string) error {
	return f.students.Remove(ctx, id)
}

// Рубрикатор

func (f *Facade) AddCriterion(ctx context.Context, in dto.AddCriterionInput) (dto.CriterionView, error) {
	return f.criteria.Add(ctx, in)
}

func (f *Facade) ListCriteria(ctx context.Context) ([]dto.CriterionView, error) {
	return f.criteria.List(ctx)
}

func (f *Facade) RemoveCriterion(ctx context.Context, id string) error {
	return f.criteria.Remove(ctx, id)
}

// Оценки

func (f *Facade) SetEvaluation(ctx context.Context, in dto.SetEvaluationInput) (dto.EvaluationView, error) {
	return f.evaluations.Set(ctx, in)
}

func (f *Facade) CorrectEvaluation(ctx context.Context, in dto.CorrectEvaluationInput) (dto.EvaluationView, error) {
	return f.evaluations.Correct(ctx, in)
}

func (f *Facade) CancelEvaluation(ctx context.Context, id string) error {
	return f.evaluations.Cancel(ctx, id)
}

func (f *Facade) ListEvaluations(ctx context.Context) ([]dto.EvaluationView, error) {
	return f.evaluations.List(ctx)
}

// Итоги

func (f *Facade) StudentTotal(ctx context.Context, studentID string) (dto.TotalRow, error) {
	return f.totals.ForStudent(ctx, studentID)
}

func (f *Facade) AllTotals(ctx context.Context) ([]dto.TotalRow, error) {
	return f.totals.All(ctx)
}

func (f *Facade) ExportTotals(ctx context.Context) (dto.ExportResult, error) {
	return f.totals.Export(ctx)
}

// StudentReport — составной сценарий: карточка студента из трёх подсистем.
// Без фасада CLI сам вызывал бы три use case и собирал результат.
func (f *Facade) StudentReport(ctx context.Context, studentID string) (dto.StudentReport, error) {
	st, err := f.students.Get(ctx, studentID)
	if err != nil {
		return dto.StudentReport{}, err
	}
	evals, err := f.evaluations.ListByStudent(ctx, studentID)
	if err != nil {
		return dto.StudentReport{}, err
	}
	total, err := f.totals.ForStudent(ctx, studentID)
	if err != nil {
		return dto.StudentReport{}, err
	}
	return dto.StudentReport{Student: st, Evaluation: evals, Total: total}, nil
}
