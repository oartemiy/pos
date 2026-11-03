package student

import (
	"context"
	"pos/internal/domain/shared"
)

type Repository interface {
	Save(ctx context.Context, s *Student) error
	FindByID(ctx context.Context, id shared.ID) (*Student, error)
	FindByEmail(ctx context.Context, email shared.Email) (*Student, error)
	List(ctx context.Context) ([]*Student, error)
	Delete(ctx context.Context, id shared.ID) error
}
