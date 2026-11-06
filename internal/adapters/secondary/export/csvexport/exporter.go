package csvexport

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"pos/internal/application/dto"
	"pos/internal/application/port"
	"strconv"
	"strings"
)

// utf8BOM нужен, чтобы был хороший парсинг кириллицы
var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

var header = []string{"ФИО", "Email", "Итог", "Оценено критериев", "Всего критериев", "Правило итога"}

type Exporter struct {
	path string
}

var _ port.TotalsExporter = (*Exporter)(nil)

func New(path string) *Exporter {
	return &Exporter{path}
}

func (e *Exporter) Export(ctx context.Context, rows []dto.TotalRow) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	buf.Write(utf8BOM)

	w := csv.NewWriter(&buf)
	w.Comma = ';'

	if err := w.Write(header); err != nil {
		return "", fmt.Errorf("csvexport: %w", err)
	}

	for _, row := range rows {
		record := []string{row.StudentName,
			row.Email,
			formatScore(row.Value),
			strconv.Itoa(row.Graded),
			strconv.Itoa(row.Criteria),
			row.Strategy}
		if err := w.Write(record); err != nil {
			return "", fmt.Errorf("csvexport: %w", err)
		}
	}
	w.Flush()

	if err := w.Error(); err != nil {
		return "", fmt.Errorf("csvexport: %w", err)
	}

	if err := writeAtomic(e.path, buf.Bytes()); err != nil {
		return "", err
	}
	return e.path, nil
}

// formatScore округляет итог до сотых для вывода и ставит десятичную запятую.
// Сам итог в домене хранится точным — округление только на границе.
func formatScore(v float64) string {
	return strings.Replace(strconv.FormatFloat(v, 'f', 2, 64), ".", ",", 1)
}

func writeAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("csvexport: mkdir: %w", err)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("csvexport: write: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("csvexport: rename: %w", err)
	}
	return nil
}
