package cli

import (
	"context"
	"fmt"

	"pos/internal/application/dto"
)

// Студенты

func (a *App) addStudent(ctx context.Context) error {
	var in dto.AddStudentInput
	var err error
	if in.LastName, err = a.in.line("Фамилия"); err != nil {
		return err
	}
	if in.FirstName, err = a.in.line("Имя"); err != nil {
		return err
	}
	if in.Patronymic, err = a.in.line("Отчество (Enter — нет)"); err != nil {
		return err
	}
	if in.Email, err = a.in.line("Email"); err != nil {
		return err
	}

	st, err := a.facade.AddStudent(ctx, in)
	if err != nil {
		return err
	}
	a.out.success("Студент добавлен: %s <%s>", st.FullName, st.Email)
	return nil
}

func (a *App) listStudents(ctx context.Context) error {
	list, err := a.facade.ListStudents(ctx)
	if err != nil {
		return err
	}
	a.out.students(list)
	return nil
}

func (a *App) removeStudent(ctx context.Context) error {
	st, err := a.pickStudent(ctx)
	if err != nil {
		return err
	}
	if err := a.confirm(fmt.Sprintf("Удалить студента %s?", st.FullName)); err != nil {
		return err
	}
	if err := a.facade.RemoveStudent(ctx, st.ID); err != nil {
		return err
	}
	a.out.success("Студент удалён")
	return nil
}

// Рубрикатор

func (a *App) addCriterion(ctx context.Context) error {
	var in dto.AddCriterionInput
	var err error
	if in.Name, err = a.in.line("Название критерия"); err != nil {
		return err
	}
	if in.MaxScore, err = a.in.number("Максимальный балл"); err != nil {
		return err
	}
	if in.Weight, err = a.in.number("Вес в итоге"); err != nil {
		return err
	}

	c, err := a.facade.AddCriterion(ctx, in)
	if err != nil {
		return err
	}
	a.out.success("Критерий «%s» добавлен: максимум %s, вес %s", c.Name, num(c.MaxScore), num(c.Weight))
	return nil
}

func (a *App) listCriteria(ctx context.Context) error {
	list, err := a.facade.ListCriteria(ctx)
	if err != nil {
		return err
	}
	a.out.criteria(list)
	return nil
}

func (a *App) removeCriterion(ctx context.Context) error {
	c, err := a.pickCriterion(ctx)
	if err != nil {
		return err
	}
	if err := a.confirm(fmt.Sprintf("Удалить критерий «%s»?", c.Name)); err != nil {
		return err
	}
	if err := a.facade.RemoveCriterion(ctx, c.ID); err != nil {
		return err
	}
	a.out.success("Критерий удалён")
	return nil
}

// Оценки

func (a *App) setEvaluation(ctx context.Context) error {
	st, err := a.pickStudent(ctx)
	if err != nil {
		return err
	}
	c, err := a.pickCriterion(ctx)
	if err != nil {
		return err
	}
	score, err := a.in.number(fmt.Sprintf("Балл (0–%s)", num(c.MaxScore)))
	if err != nil {
		return err
	}
	comment, err := a.in.line("Комментарий (Enter — без комментария)")
	if err != nil {
		return err
	}

	e, err := a.facade.SetEvaluation(ctx, dto.SetEvaluationInput{
		StudentID:   st.ID,
		CriterionID: c.ID,
		Score:       score,
		Comment:     comment,
	})
	if err != nil {
		return err
	}
	a.out.success("Оценка выставлена: %s — %s: %s из %s",
		e.StudentName, e.CriterionName, num(e.Score), num(e.MaxScore))
	return nil
}

func (a *App) listEvaluations(ctx context.Context) error {
	list, err := a.facade.ListEvaluations(ctx)
	if err != nil {
		return err
	}
	a.out.evaluations(list)
	return nil
}

func (a *App) correctEvaluation(ctx context.Context) error {
	e, err := a.pickEvaluation(ctx)
	if err != nil {
		return err
	}
	score, err := a.in.number(fmt.Sprintf("Новый балл (0–%s)", num(e.MaxScore)))
	if err != nil {
		return err
	}
	comment, err := a.in.line("Комментарий (Enter — без комментария)")
	if err != nil {
		return err
	}

	updated, err := a.facade.CorrectEvaluation(ctx, dto.CorrectEvaluationInput{
		EvaluationID: e.ID,
		Score:        score,
		Comment:      comment,
	})
	if err != nil {
		return err
	}
	a.out.success("Оценка исправлена: %s — %s: %s из %s",
		updated.StudentName, updated.CriterionName, num(updated.Score), num(updated.MaxScore))
	return nil
}

func (a *App) cancelEvaluation(ctx context.Context) error {
	e, err := a.pickEvaluation(ctx)
	if err != nil {
		return err
	}
	if err := a.confirm(fmt.Sprintf("Отменить оценку %s — %s?", e.StudentName, e.CriterionName)); err != nil {
		return err
	}
	if err := a.facade.CancelEvaluation(ctx, e.ID); err != nil {
		return err
	}
	a.out.success("Оценка отменена")
	return nil
}

// Итоги

func (a *App) studentReport(ctx context.Context) error {
	st, err := a.pickStudent(ctx)
	if err != nil {
		return err
	}
	r, err := a.facade.StudentReport(ctx, st.ID)
	if err != nil {
		return err
	}
	a.out.report(r)
	return nil
}

func (a *App) allTotals(ctx context.Context) error {
	rows, err := a.facade.AllTotals(ctx)
	if err != nil {
		return err
	}
	a.out.totals(rows)
	return nil
}

func (a *App) exportTotals(ctx context.Context) error {
	res, err := a.facade.ExportTotals(ctx)
	if err != nil {
		return err
	}
	a.out.success("Выгружено строк: %d → %s", res.Rows, res.Location)
	return nil
}

// Выбор из списка. Пользователь вводит номер строки, а не UUID.

func (a *App) pickStudent(ctx context.Context) (dto.StudentView, error) {
	list, err := a.facade.ListStudents(ctx)
	if err != nil {
		return dto.StudentView{}, err
	}
	if len(list) == 0 {
		return dto.StudentView{}, errNoStudents
	}
	return choose(a.in, list, a.out.students, "Номер студента")
}

func (a *App) pickCriterion(ctx context.Context) (dto.CriterionView, error) {
	list, err := a.facade.ListCriteria(ctx)
	if err != nil {
		return dto.CriterionView{}, err
	}
	if len(list) == 0 {
		return dto.CriterionView{}, errNoCriteria
	}
	return choose(a.in, list, a.out.criteria, "Номер критерия")
}

func (a *App) pickEvaluation(ctx context.Context) (dto.EvaluationView, error) {
	list, err := a.facade.ListEvaluations(ctx)
	if err != nil {
		return dto.EvaluationView{}, err
	}
	if len(list) == 0 {
		return dto.EvaluationView{}, errNoEvaluations
	}
	return choose(a.in, list, a.out.evaluations, "Номер оценки")
}

func choose[T any](in *input, list []T, show func([]T), prompt string) (T, error) {
	var zero T
	show(list)
	i, err := in.index(prompt, len(list))
	if err != nil {
		return zero, err
	}
	return list[i], nil
}

func (a *App) confirm(prompt string) error {
	ok, err := a.in.confirm(prompt)
	if err != nil {
		return err
	}
	if !ok {
		return errCancelled
	}
	return nil
}
