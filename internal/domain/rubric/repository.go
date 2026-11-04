package rubric

import "context"

type Repository interface {
	Get(ctx context.Context) (*Rubric, error)
	Save(ctx context.Context, r *Rubric) error
}
