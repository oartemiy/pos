package usecase

import (
	"context"
	"pos/internal/application/dto"
)

func (uc *Totals) Export(ctx context.Context) (dto.ExportResult, error) {
	rows, err := uc.All(ctx)
	if err != nil {
		return dto.ExportResult{}, err
	}
	location, err := uc.exporter.Export(ctx, rows)
	if err != nil {
		return dto.ExportResult{}, err
	}
	return dto.ExportResult{Location: location, Rows: len(rows)}, nil
}
