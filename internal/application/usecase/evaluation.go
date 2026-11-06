package usecase

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"pos/internal/application/dto"
	"pos/internal/domain/evaluation"
	"pos/internal/domain/rubric"
	"pos/internal/domain/shared"
	"pos/internal/domain/student"
)

type Evaluations struct {
	students    student.Repository
	rubrics     rubric.Repository
	evaluations evaluation.Repository
	factory     *evaluation.Factory
	now         func() time.Time
}

func NewEvaluations(
	s student.Repository,
	r rubric.Repository,
	e evaluation.Repository,
	f *evaluation.Factory,
	now func() time.Time,
) *Evaluations {
	return &Evaluations{students: s, rubrics: r, evaluations: e, factory: f, now: now}
}

// Set выставляет оценку. Инварианты оценки проверяет фабрика,
// а здесь — то, что требует других агрегатов.
func (uc *Evaluations) Set(ctx context.Context, in dto.SetEvaluationInput) (dto.EvaluationView, error) {
	studentID, err := shared.ParseID(in.StudentID)
	if err != nil {
		return dto.EvaluationView{}, err
	}
	criterionID, err := shared.ParseID(in.CriterionID)
	if err != nil {
		return dto.EvaluationView{}, err
	}

	st, err := uc.students.FindByID(ctx, studentID)
	if err != nil {
		return dto.EvaluationView{}, notFoundAs(err, shared.ErrStudentNotFound)
	}
	rb, err := loadRubric(ctx, uc.rubrics)
	if err != nil {
		return dto.EvaluationView{}, err
	}

	// Одна оценка на пару «студент + критерий»: повторное выставление — это исправление.
	_, err = uc.evaluations.FindByStudentAndCriterion(ctx, studentID, criterionID)
	switch {
	case err == nil:
		return dto.EvaluationView{}, shared.ErrEvaluationExists
	case !errors.Is(err, shared.ErrNotFound):
		return dto.EvaluationView{}, err
	}

	e, err := uc.factory.Create(rb, studentID, criterionID, in.Score, in.Comment)
	if err != nil {
		return dto.EvaluationView{}, err
	}
	if err := uc.evaluations.Save(ctx, e); err != nil {
		return dto.EvaluationView{}, err
	}

	c, _ := rb.Criterion(criterionID) // существование уже проверила фабрика
	return toEvaluationView(e, st, c), nil
}

func (uc *Evaluations) Correct(ctx context.Context, in dto.CorrectEvaluationInput) (dto.EvaluationView, error) {
	id, err := shared.ParseID(in.EvaluationID)
	if err != nil {
		return dto.EvaluationView{}, err
	}
	e, err := uc.evaluations.FindByID(ctx, id)
	if err != nil {
		return dto.EvaluationView{}, notFoundAs(err, shared.ErrEvaluationNotFound)
	}
	rb, err := loadRubric(ctx, uc.rubrics)
	if err != nil {
		return dto.EvaluationView{}, err
	}
	c, err := rb.Criterion(e.CriterionID())
	if err != nil {
		return dto.EvaluationView{}, err
	}

	score, err := evaluation.NewScore(in.Score)
	if err != nil {
		return dto.EvaluationView{}, err
	}
	comment, err := evaluation.NewComment(in.Comment)
	if err != nil {
		return dto.EvaluationView{}, err
	}
	if err := e.Correct(c, score, comment, uc.now()); err != nil {
		return dto.EvaluationView{}, err
	}
	if err := uc.evaluations.Save(ctx, e); err != nil {
		return dto.EvaluationView{}, err
	}

	st, err := uc.students.FindByID(ctx, e.StudentID())
	if err != nil {
		return dto.EvaluationView{}, notFoundAs(err, shared.ErrStudentNotFound)
	}
	return toEvaluationView(e, st, c), nil
}

func (uc *Evaluations) Cancel(ctx context.Context, rawID string) error {
	id, err := shared.ParseID(rawID)
	if err != nil {
		return err
	}
	return notFoundAs(uc.evaluations.Delete(ctx, id), shared.ErrEvaluationNotFound)
}

func (uc *Evaluations) List(ctx context.Context) ([]dto.EvaluationView, error) {
	evals, err := uc.evaluations.List(ctx)
	if err != nil {
		return nil, err
	}
	return uc.views(ctx, evals)
}

func (uc *Evaluations) ListByStudent(ctx context.Context, rawStudentID string) ([]dto.EvaluationView, error) {
	id, err := shared.ParseID(rawStudentID)
	if err != nil {
		return nil, err
	}
	evals, err := uc.evaluations.FindByStudent(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.views(ctx, evals)
}

// views дополняет оценки именами студентов и критериев для вывода.
func (uc *Evaluations) views(ctx context.Context, evals []*evaluation.Evaluation) ([]dto.EvaluationView, error) {
	if len(evals) == 0 {
		return []dto.EvaluationView{}, nil
	}
	students, err := uc.students.List(ctx)
	if err != nil {
		return nil, err
	}
	byID := make(map[shared.ID]*student.Student, len(students))
	for _, s := range students {
		byID[s.ID()] = s
	}
	rb, err := loadRubric(ctx, uc.rubrics)
	if err != nil {
		return nil, err
	}

	views := make([]dto.EvaluationView, 0, len(evals))
	for _, e := range evals {
		st, ok := byID[e.StudentID()]
		if !ok {
			return nil, fmt.Errorf("оценка %s: %w", e.ID(), shared.ErrStudentNotFound)
		}
		c, err := rb.Criterion(e.CriterionID())
		if err != nil {
			return nil, fmt.Errorf("оценка %s: %w", e.ID(), err)
		}
		views = append(views, toEvaluationView(e, st, c))
	}

	slices.SortFunc(views, func(a, b dto.EvaluationView) int {
		return cmp.Or(
			cmp.Compare(a.StudentName, b.StudentName),
			cmp.Compare(a.CriterionName, b.CriterionName),
		)
	})
	return views, nil
}

func toEvaluationView(e *evaluation.Evaluation, st *student.Student, c rubric.Criterion) dto.EvaluationView {
	return dto.EvaluationView{
		ID:            e.ID().String(),
		StudentID:     st.ID().String(),
		StudentName:   st.Name().Short(),
		CriterionID:   c.ID().String(),
		CriterionName: c.Name(),
		Score:         e.Score().Value(),
		MaxScore:      c.MaxScore().Value(),
		Comment:       e.Comment().String(),
		GradedAt:      e.GradedAt(),
	}
}
