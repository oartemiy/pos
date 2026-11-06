package usecase

import (
	"context"

	"pos/internal/application/dto"
	"pos/internal/application/port"
	"pos/internal/domain/evaluation"
	"pos/internal/domain/grading"
	"pos/internal/domain/rubric"
	"pos/internal/domain/shared"
	"pos/internal/domain/student"
)

// Totals — сценарии итога. Итог не хранится: он каждый раз считается
// Domain Service из текущих оценок и рубрикатора.
type Totals struct {
	students    student.Repository
	rubrics     rubric.Repository
	evaluations evaluation.Repository
	grading     *grading.Service
	exporter    port.TotalsExporter
}

func NewTotals(
	s student.Repository,
	r rubric.Repository,
	e evaluation.Repository,
	g *grading.Service,
	exp port.TotalsExporter,
) *Totals {
	return &Totals{students: s, rubrics: r, evaluations: e, grading: g, exporter: exp}
}

func (uc *Totals) ForStudent(ctx context.Context, rawID string) (dto.TotalRow, error) {
	id, err := shared.ParseID(rawID)
	if err != nil {
		return dto.TotalRow{}, err
	}
	st, err := uc.students.FindByID(ctx, id)
	if err != nil {
		return dto.TotalRow{}, notFoundAs(err, shared.ErrStudentNotFound)
	}
	rb, err := loadRubric(ctx, uc.rubrics)
	if err != nil {
		return dto.TotalRow{}, err
	}
	evals, err := uc.evaluations.FindByStudent(ctx, id)
	if err != nil {
		return dto.TotalRow{}, err
	}

	total, err := uc.grading.Calculate(id, rb, evals)
	if err != nil {
		return dto.TotalRow{}, err
	}
	return toTotalRow(st, total), nil
}

func (uc *Totals) All(ctx context.Context) ([]dto.TotalRow, error) {
	students, err := uc.students.List(ctx)
	if err != nil {
		return nil, err
	}
	sortStudents(students)

	rb, err := loadRubric(ctx, uc.rubrics)
	if err != nil {
		return nil, err
	}
	evals, err := uc.evaluations.List(ctx)
	if err != nil {
		return nil, err
	}
	byStudent := make(map[shared.ID][]*evaluation.Evaluation)
	for _, e := range evals {
		byStudent[e.StudentID()] = append(byStudent[e.StudentID()], e)
	}

	rows := make([]dto.TotalRow, 0, len(students))
	for _, st := range students {
		total, err := uc.grading.Calculate(st.ID(), rb, byStudent[st.ID()])
		if err != nil {
			return nil, err
		}
		rows = append(rows, toTotalRow(st, total))
	}
	return rows, nil
}

func toTotalRow(st *student.Student, t grading.Total) dto.TotalRow {
	return dto.TotalRow{
		StudentID:   st.ID().String(),
		StudentName: st.Name().Short(),
		Email:       st.Email().String(),
		Value:       t.Value(),
		Graded:      t.Graded(),
		Criteria:    t.Criteria(),
		Strategy:    t.Strategy(),
	}
}
