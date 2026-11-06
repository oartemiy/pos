package usecase

import (
	"context"
	"errors"
	"pos/internal/domain/rubric"
	"pos/internal/domain/shared"
)

func notFoundAs(err, target error) error {
	if errors.Is(err, shared.ErrNotFound) {
		return target
	}
	return err
}

func loadRubric(ctx context.Context, repo rubric.Repository) (*rubric.Rubric, error) {
	rb, err := repo.Get(ctx)
	return rb, notFoundAs(err, shared.ErrEmptyRubric)
}
