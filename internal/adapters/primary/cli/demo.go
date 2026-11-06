package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"pos/internal/application/dto"
	"pos/internal/domain/shared"
)

// runDemo прогоняет сквозной сценарий из ТЗ:
// студенты → критерии → оценки → нарушения правил → итоги → экспорт.
// Повторный запуск не падает: существующие данные переиспользуются.
func (a *App) runDemo(ctx context.Context) error {
	a.out.info("Шаг 1. Заводим студентов")
	studentInputs := []dto.AddStudentInput{
		{LastName: "Иванов", FirstName: "Иван", Patronumic: "Петрович", Email: "ivan.ivanov@edu.hse.ru"},
		{LastName: "Петрова", FirstName: "Анна", Patronumic: "Сергеевна", Email: "anna.petrova@edu.hse.ru"},
		{LastName: "Smith", FirstName: "John", Email: "john.smith@edu.hse.ru"},
	}
	students := make([]dto.StudentView, 0, len(studentInputs))
	for _, in := range studentInputs {
		st, err := a.ensureStudent(ctx, in)
		if err != nil {
			return err
		}
		a.out.success("%s <%s>", st.FullName, st.Email)
		students = append(students, st)
	}

	a.out.info("Шаг 2. Заполняем рубрикатор")
	criterionInputs := []dto.AddCriterionInput{
		{Name: "Архитектура", MaxScore: 10, Weight: 2},
		{Name: "Паттерны GoF", MaxScore: 10, Weight: 2},
		{Name: "Тесты", MaxScore: 5, Weight: 1},
	}
	criteria := make([]dto.CriterionView, 0, len(criterionInputs))
	for _, in := range criterionInputs {
		c, err := a.ensureCriterion(ctx, in)
		if err != nil {
			return err
		}
		a.out.success("«%s»: максимум %s, вес %s", c.Name, num(c.MaxWeight), num(c.Weight))
		criteria = append(criteria, c)
	}

	a.out.info("Шаг 3. Выставляем оценки")
	grades := []struct {
		student, criterion int
		score              float64
		comment            string
	}{
		{0, 0, 9, "чистые слои"},
		{0, 1, 8, ""},
		{0, 2, 5, ""},
		{1, 0, 7, ""},
		{1, 1, 10, "паттерны по делу"},
		{2, 0, 6, ""}, // у третьего студента оценены не все критерии
	}
	var first dto.EvaluationView
	for i, g := range grades {
		e, err := a.grade(ctx, students[g.student].ID, criteria[g.criterion].ID, g.score, g.comment)
		if err != nil {
			return err
		}
		if i == 0 {
			first = e
		}
		a.out.success("%s — %s: %s из %s", e.StudentName, e.CriterionName, num(e.Score), num(e.MaxScore))
	}

	a.out.info("Шаг 4. Пробуем нарушить правила домена")
	_, err := a.facade.CorrectEvaluation(ctx, dto.CorrectEvaluationInput{EvaluationID: first.ID, Score: 12})
	if err := a.expectRejected("исправить оценку на 12 из 10", err, shared.ErrScoreOutOfRange); err != nil {
		return err
	}
	_, err = a.facade.SetEvaluation(ctx, dto.SetEvaluationInput{
		StudentID:   students[2].ID,
		CriterionID: "00000000-0000-0000-0000-000000000000",
		Score:       5,
	})
	if err := a.expectRejected("оценить критерий, которого нет в рубрикаторе", err, shared.ErrCriterionNotFound); err != nil {
		return err
	}

	a.out.info("Шаг 5. Считаем итоги")
	rows, err := a.facade.AllTotals(ctx)
	if err != nil {
		return err
	}
	a.out.totals(rows)

	a.out.info("Шаг 6. Экспортируем итоги")
	res, err := a.facade.ExportTotals(ctx)
	if err != nil {
		return err
	}
	a.out.success("Выгружено строк: %d → %s", res.Rows, res.Location)
	return nil
}

func (a *App) ensureStudent(ctx context.Context, in dto.AddStudentInput) (dto.StudentView, error) {
	st, err := a.facade.AddStudent(ctx, in)
	if !errors.Is(err, shared.ErrEmailTaken) {
		return st, err
	}
	list, err := a.facade.ListStudents(ctx)
	if err != nil {
		return dto.StudentView{}, err
	}
	for _, s := range list {
		if strings.EqualFold(s.Email, in.Email) {
			return s, nil
		}
	}
	return dto.StudentView{}, shared.ErrStudentNotFound
}

func (a *App) ensureCriterion(ctx context.Context, in dto.AddCriterionInput) (dto.CriterionView, error) {
	c, err := a.facade.AddCriterion(ctx, in)
	if !errors.Is(err, shared.ErrDuplicateCriterion) {
		return c, err
	}
	list, err := a.facade.ListCriteria(ctx)
	if err != nil {
		return dto.CriterionView{}, err
	}
	for _, existing := range list {
		if strings.EqualFold(existing.Name, in.Name) {
			return existing, nil
		}
	}
	return dto.CriterionView{}, shared.ErrCriterionNotFound
}

// grade выставляет оценку, а если она уже есть — исправляет её.
func (a *App) grade(ctx context.Context, studentID, criterionID string, score float64, comment string) (dto.EvaluationView, error) {
	e, err := a.facade.SetEvaluation(ctx, dto.SetEvaluationInput{
		StudentID:   studentID,
		CriterionID: criterionID,
		Score:       score,
		Comment:     comment,
	})
	if !errors.Is(err, shared.ErrEvaluationExists) {
		return e, err
	}
	list, err := a.facade.ListEvaluations(ctx)
	if err != nil {
		return dto.EvaluationView{}, err
	}
	for _, existing := range list {
		if existing.StudentID == studentID && existing.CriterionID == criterionID {
			return a.facade.CorrectEvaluation(ctx, dto.CorrectEvaluationInput{
				EvaluationID: existing.ID,
				Score:        score,
				Comment:      comment,
			})
		}
	}
	return dto.EvaluationView{}, shared.ErrEvaluationNotFound
}

func (a *App) expectRejected(action string, err, want error) error {
	switch {
	case errors.Is(err, want):
		a.out.success("Нельзя %s: %v", action, err)
		return nil
	case err == nil:
		return fmt.Errorf("демо: удалось %s, хотя правила это запрещают", action)
	default:
		return err
	}
}
