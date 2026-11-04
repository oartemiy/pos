package evaluation

import (
	"context"

	"pos/internal/domain/shared"
)

type Repository interface {
	Save(ctx context.Context, e *Evaluation) error
	FindByID(ctx context.Context, id shared.ID) (*Evaluation, error)
	FindByStudentAndCriterion(ctx context.Context, studentID, criterionID shared.ID) (*Evaluation, error)
	FindByStudent(ctx context.Context, studentID shared.ID) ([]*Evaluation, error)
	FindByCriterion(ctx context.Context, criterionID shared.ID) ([]*Evaluation, error)
	List(ctx context.Context) ([]*Evaluation, error)
	Delete(ctx context.Context, id shared.ID) error
}
