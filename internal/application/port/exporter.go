package port

import (
	"context"
	"pos/internal/application/dto"
)

type TotalsExporter interface {
	Export(ctx context.Context, rows []dto.TotalRow) (string, error)
}