package di

import (
	"fmt"
	"path/filepath"
	"time"

	"go.uber.org/dig"

	"pos/internal/adapters/secondary/export/csvexport"
	"pos/internal/adapters/secondary/persistence/jsonfile"
	"pos/internal/application/facade"
	"pos/internal/application/port"
	"pos/internal/application/usecase"
	"pos/internal/domain/evaluation"
	"pos/internal/domain/grading"
	"pos/internal/domain/rubric"
	"pos/internal/domain/student"
)

type Config struct {
	DataDir     string
	CourseTitle string
	Strategy    string
}

func Build(cfg Config) (*dig.Container, error) {
	c := dig.New()

	providers := []any{
		func() Config { return cfg },

		// Secondary-адаптеры: реализации портов домена и приложения.
		func(cfg Config) student.Repository {
			return jsonfile.NewStudentRepo(filepath.Join(cfg.DataDir, "students.json"))
		},
		func(cfg Config) rubric.Repository {
			return jsonfile.NewRubricRepo(filepath.Join(cfg.DataDir, "rubric.json"))
		},
		func(cfg Config) evaluation.Repository {
			return jsonfile.NewEvaluationRepo(filepath.Join(cfg.DataDir, "evaluations.json"))
		},
		func(cfg Config) port.TotalsExporter {
			return csvexport.New(filepath.Join(cfg.DataDir, "totals.csv"))
		},

		newStrategy,
		grading.NewService,
		func() *evaluation.Factory { return evaluation.NewFactory(time.Now) },

		usecase.NewStudents,
		func(r rubric.Repository, e evaluation.Repository, cfg Config) *usecase.Criteria {
			return usecase.NewCriteria(r, e, cfg.CourseTitle)
		},
		func(s student.Repository, r rubric.Repository, e evaluation.Repository, f *evaluation.Factory) *usecase.Evaluations {
			return usecase.NewEvaluations(s, r, e, f, time.Now)
		},
		usecase.NewTotals,
		facade.New,
	}

	for _, p := range providers {
		if err := c.Provide(p); err != nil {
			return nil, fmt.Errorf("di: %w", err)
		}
	}
	return c, nil
}

func newStrategy(cfg Config) (grading.Strategy, error) {
	switch cfg.Strategy {
	case "", "average":
		return grading.WeightedAverage{}, nil
	case "sum":
		return grading.WeightedSum{}, nil
	default:
		return nil, fmt.Errorf("неизвестное правило итога %q (доступны: average, sum)", cfg.Strategy)
	}
}
