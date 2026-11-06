package usecase

import (
	"context"
	"errors"

	"pos/internal/application/dto"
	"pos/internal/domain/evaluation"
	"pos/internal/domain/rubric"
	"pos/internal/domain/shared"
)

type Criteria struct {
	rubrics     rubric.Repository
	evaluations evaluation.Repository
	courseTitle string
}

// courseTitle — название рубрикатора, который создаётся при добавлении первого критерия.
func NewCriteria(r rubric.Repository, e evaluation.Repository, courseTitle string) *Criteria {
	return &Criteria{rubrics: r, evaluations: e, courseTitle: courseTitle}
}

func (uc *Criteria) Add(ctx context.Context, in dto.AddCriterionInput) (dto.CriterionView, error) {
	max, err := rubric.NewMaxScore(in.MaxScore)
	if err != nil {
		return dto.CriterionView{}, err
	}
	w, err := rubric.NewWeight(in.Weight)
	if err != nil {
		return dto.CriterionView{}, err
	}

	rb, err := uc.rubrics.Get(ctx)
	if errors.Is(err, shared.ErrNotFound) {
		rb, err = rubric.New(uc.courseTitle)
	}
	if err != nil {
		return dto.CriterionView{}, err
	}

	c, err := rb.AddCriterion(in.Name, max, w)
	if err != nil {
		return dto.CriterionView{}, err
	}
	if err := uc.rubrics.Save(ctx, rb); err != nil {
		return dto.CriterionView{}, err
	}
	return toCriterionView(c), nil
}

func (uc *Criteria) List(ctx context.Context) ([]dto.CriterionView, error) {
	rb, err := uc.rubrics.Get(ctx)
	if errors.Is(err, shared.ErrNotFound) {
		return []dto.CriterionView{}, nil
	}
	if err != nil {
		return nil, err
	}

	criteria := rb.Criteria()
	views := make([]dto.CriterionView, 0, len(criteria))
	for _, c := range criteria {
		views = append(views, toCriterionView(c))
	}
	return views, nil
}

// Remove удаляет критерий, только если по нему нет оценок.
// Правило между агрегатами: рубрикатор об оценках не знает.
func (uc *Criteria) Remove(ctx context.Context, rawID string) error {
	id, err := shared.ParseID(rawID)
	if err != nil {
		return err
	}
	rb, err := uc.rubrics.Get(ctx)
	if err != nil {
		return notFoundAs(err, shared.ErrCriterionNotFound)
	}
	if _, err := rb.Criterion(id); err != nil {
		return err
	}

	evals, err := uc.evaluations.FindByCriterion(ctx, id)
	if err != nil {
		return err
	}
	if len(evals) > 0 {
		return shared.ErrCriterionHasEvaluations
	}

	if err := rb.RemoveCriterion(id); err != nil {
		return err
	}
	return uc.rubrics.Save(ctx, rb)
}

func toCriterionView(c rubric.Criterion) dto.CriterionView {
	return dto.CriterionView{
		ID:        c.ID().String(),
		Name:      c.Name(),
		MaxScore: c.MaxScore().Value(),
		Weight:    c.Weight().Value(),
	}
}
