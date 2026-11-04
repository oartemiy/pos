package jsonfile

import (
	"context"

	"pos/internal/domain/rubric"
	"pos/internal/domain/shared"
)

type RubricRepo struct {
	store *store[rubricRecord]
}

// To check intercace implementation
var _ rubric.Repository = (*RubricRepo)(nil)

func NewRubricRepo(path string) *RubricRepo {
	return &RubricRepo{store: newStore[rubricRecord](path)}
}

func (r *RubricRepo) Get(ctx context.Context) (*rubric.Rubric, error) {
	items, err := r.store.read(ctx)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, shared.ErrNotFound
	}
	return toRubric(items[0])
}

// Save сохраняет агрегат целиком: состав критериев всегда согласован.
func (r *RubricRepo) Save(ctx context.Context, rb *rubric.Rubric) error {
	rec := toRubricRecord(rb)
	return r.store.update(ctx, func(_ []rubricRecord) ([]rubricRecord, error) {
		return []rubricRecord{rec}, nil // один курс — один рубрикатор
	})
}
