package jsonfile

import (
	"context"

	"pos/internal/domain/evaluation"
	"pos/internal/domain/shared"
)

type EvaluationRepo struct {
	store *store[evaluationRecord]
}

var _ evaluation.Repository = (*EvaluationRepo)(nil)

func NewEvaluationRepo(path string) *EvaluationRepo {
	return &EvaluationRepo{store: newStore[evaluationRecord](path)}
}

func (r *EvaluationRepo) Save(ctx context.Context, e *evaluation.Evaluation) error {
	rec := toEvaluationRecord(e)
	return r.store.update(ctx, func(items []evaluationRecord) ([]evaluationRecord, error) {
		for i := range items {
			if items[i].ID == rec.ID {
				items[i] = rec // исправление существующей оценки
				return items, nil
			}
		}
		return append(items, rec), nil
	})
}

func (r *EvaluationRepo) FindByID(ctx context.Context, id shared.ID) (*evaluation.Evaluation, error) {
	return r.findOne(ctx, func(rec evaluationRecord) bool {
		return rec.ID == id.String()
	})
}

func (r *EvaluationRepo) FindByStudentAndCriterion(ctx context.Context, studentID, criterionID shared.ID) (*evaluation.Evaluation, error) {
	sid, cid := studentID.String(), criterionID.String()
	return r.findOne(ctx, func(rec evaluationRecord) bool {
		return rec.StudentID == sid && rec.CriterionID == cid
	})
}

func (r *EvaluationRepo) FindByStudent(ctx context.Context, studentID shared.ID) ([]*evaluation.Evaluation, error) {
	sid := studentID.String()
	return r.findAll(ctx, func(rec evaluationRecord) bool { return rec.StudentID == sid })
}

func (r *EvaluationRepo) FindByCriterion(ctx context.Context, criterionID shared.ID) ([]*evaluation.Evaluation, error) {
	cid := criterionID.String()
	return r.findAll(ctx, func(rec evaluationRecord) bool { return rec.CriterionID == cid })
}

func (r *EvaluationRepo) List(ctx context.Context) ([]*evaluation.Evaluation, error) {
	return r.findAll(ctx, func(evaluationRecord) bool { return true })
}

func (r *EvaluationRepo) Delete(ctx context.Context, id shared.ID) error {
	return r.store.update(ctx, func(items []evaluationRecord) ([]evaluationRecord, error) {
		for i := range items {
			if items[i].ID == id.String() {
				return append(items[:i], items[i+1:]...), nil
			}
		}
		return nil, shared.ErrNotFound
	})
}

func (r *EvaluationRepo) findOne(ctx context.Context, match func(evaluationRecord) bool) (*evaluation.Evaluation, error) {
	items, err := r.store.read(ctx)
	if err != nil {
		return nil, err
	}
	for _, rec := range items {
		if match(rec) {
			return toEvaluation(rec)
		}
	}
	return nil, shared.ErrNotFound
}

func (r *EvaluationRepo) findAll(ctx context.Context, match func(evaluationRecord) bool) ([]*evaluation.Evaluation, error) {
	items, err := r.store.read(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*evaluation.Evaluation, 0)
	for _, rec := range items {
		if !match(rec) {
			continue
		}
		e, err := toEvaluation(rec)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, nil
}
