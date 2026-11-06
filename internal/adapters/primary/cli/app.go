// Package cli — консольный primary-адаптер ПОС. Работает только через фасад
// приложения и ничего не знает о хранилище и доменных объектах.
package cli

import (
	"context"
	"errors"
	"io"

	"pos/internal/application/facade"
)

var (
	errUnknownCommand = errors.New("нет такого пункта меню")
	errCancelled      = errors.New("действие отменено")
	errNoStudents     = errors.New("студентов пока нет, сначала добавьте их")
	errNoCriteria     = errors.New("критериев пока нет, сначала заполните рубрикатор")
	errNoEvaluations  = errors.New("оценок пока нет")
)

type command struct {
	key     string
	section string
	title   string
	run     func(ctx context.Context) error
}

type App struct {
	facade   *facade.Facade
	in       *input
	out      *output
	commands []command
}

func New(f *facade.Facade, r io.Reader, w io.Writer) *App {
	a := &App{facade: f, in: newInput(r, w), out: newOutput(w)}
	a.commands = []command{
		{"1", "Студенты", "Добавить студента", a.addStudent},
		{"2", "Студенты", "Список студентов", a.listStudents},
		{"3", "Студенты", "Удалить студента", a.removeStudent},
		{"4", "Рубрикатор", "Добавить критерий", a.addCriterion},
		{"5", "Рубрикатор", "Список критериев", a.listCriteria},
		{"6", "Рубрикатор", "Удалить критерий", a.removeCriterion},
		{"7", "Оценки", "Выставить оценку", a.setEvaluation},
		{"8", "Оценки", "Список оценок", a.listEvaluations},
		{"9", "Оценки", "Исправить оценку", a.correctEvaluation},
		{"10", "Оценки", "Отменить оценку", a.cancelEvaluation},
		{"11", "Итоги", "Карточка студента", a.studentReport},
		{"12", "Итоги", "Итоги всех студентов", a.allTotals},
		{"13", "Итоги", "Экспорт итогов в CSV", a.exportTotals},
		{"14", "Демо", "Сквозной демо-сценарий", a.runDemo},
	}
	return a
}

// Run — главный цикл. Ошибки сценариев (например, «балл выходит за пределы
// критерия») выводятся пользователю, и программа продолжает работу.
func (a *App) Run(ctx context.Context) error {
	a.in.bind(ctx)
	a.out.title("ПОС — преподавательская оценка студентов")
	for {
		if err := ctx.Err(); err != nil {
			return a.stop(err)
		}
		a.out.menu(a.commands)
		choice, err := a.in.line("Пункт меню")
		if err != nil {
			return a.stop(err)
		}
		if choice == "0" {
			return a.stop(io.EOF)
		}

		cmd, ok := a.find(choice)
		if !ok {
			a.out.failure(errUnknownCommand)
			continue
		}

		a.out.section(cmd.title)
		switch err := cmd.run(ctx); {
		case err == nil:
		case errors.Is(err, io.EOF), ctx.Err() != nil:
			return a.stop(err)
		case errors.Is(err, errCancelled):
			a.out.info("Действие отменено")
		default:
			a.out.failure(err)
		}
	}
}

func (a *App) find(key string) (command, bool) {
	for _, c := range a.commands {
		if c.key == key {
			return c, true
		}
	}
	return command{}, false
}

// stop завершает работу: конец ввода (пункт «0» или Ctrl+D) и Ctrl+C —
// штатный выход, любая другая ошибка ввода возвращается в main.
func (a *App) stop(err error) error {
	if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
		a.out.println("\nДо встречи!")
		return nil
	}
	return err
}
