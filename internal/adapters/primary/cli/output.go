package cli

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"
	"unicode/utf8"

	"pos/internal/application/dto"
)

type output struct {
	w io.Writer
}

func newOutput(w io.Writer) *output {
	return &output{w: w}
}

func (o *output) println(a ...any) { fmt.Fprintln(o.w, a...) }

func (o *output) title(s string) {
	fmt.Fprintf(o.w, "\n%s\n%s\n", s, strings.Repeat("═", utf8.RuneCountInString(s)))
}

func (o *output) section(s string) { fmt.Fprintf(o.w, "\n— %s —\n", s) }

func (o *output) info(format string, args ...any) {
	fmt.Fprintf(o.w, "  "+format+"\n", args...)
}

func (o *output) success(format string, args ...any) {
	fmt.Fprintf(o.w, "✓ "+format+"\n", args...)
}

func (o *output) failure(err error) { fmt.Fprintf(o.w, "✗ %v\n", err) }

func (o *output) menu(cmds []command) {
	section := ""
	for _, c := range cmds {
		if c.section != section {
			section = c.section
			fmt.Fprintf(o.w, "\n%s\n", section)
		}
		fmt.Fprintf(o.w, "  %2s. %s\n", c.key, c.title)
	}
	fmt.Fprintf(o.w, "\n   0. Выход\n\n")
}

func (o *output) table(headers []string, rows [][]string) {
	tw := tabwriter.NewWriter(o.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, strings.Join(headers, "\t"))
	for _, r := range rows {
		fmt.Fprintln(tw, strings.Join(r, "\t"))
	}
	tw.Flush()
}

func (o *output) students(list []dto.StudentView) {
	if len(list) == 0 {
		o.info("Студентов пока нет")
		return
	}
	rows := make([][]string, 0, len(list))
	for i, s := range list {
		rows = append(rows, []string{strconv.Itoa(i + 1), s.FullName, s.Email})
	}
	o.table([]string{"№", "ФИО", "Email"}, rows)
}

func (o *output) criteria(list []dto.CriterionView) {
	if len(list) == 0 {
		o.info("Рубрикатор пуст")
		return
	}
	rows := make([][]string, 0, len(list))
	for i, c := range list {
		rows = append(rows, []string{strconv.Itoa(i + 1), c.Name, num(c.MaxScore), num(c.Weight)})
	}
	o.table([]string{"№", "Критерий", "Максимум", "Вес"}, rows)
}

func (o *output) evaluations(list []dto.EvaluationView) {
	if len(list) == 0 {
		o.info("Оценок пока нет")
		return
	}
	rows := make([][]string, 0, len(list))
	for i, e := range list {
		comment := e.Comment
		if comment == "" {
			comment = "—"
		}
		rows = append(rows, []string{
			strconv.Itoa(i + 1),
			e.StudentName,
			e.CriterionName,
			num(e.Score) + " / " + num(e.MaxScore),
			e.GradedAt.Local().Format("02.01.2006 15:04"),
			comment,
		})
	}
	o.table([]string{"№", "Студент", "Критерий", "Балл", "Дата", "Комментарий"}, rows)
}

func (o *output) totals(list []dto.TotalRow) {
	if len(list) == 0 {
		o.info("Студентов пока нет")
		return
	}
	rows := make([][]string, 0, len(list))
	for i, t := range list {
		rows = append(rows, []string{
			strconv.Itoa(i + 1),
			t.StudentName,
			total(t.Value),
			fmt.Sprintf("%d из %d", t.Graded, t.Criteria),
		})
	}
	o.table([]string{"№", "Студент", "Итог", "Оценено критериев"}, rows)
	o.info("Правило итога: %s", list[0].Strategy)
}

func (o *output) report(r dto.StudentReport) {
	o.println(r.Student.FullName + " <" + r.Student.Email + ">")
	o.evaluations(r.Evaluation)
	o.info("Итог: %s (оценено %d из %d, %s)",
		total(r.Total.Value), r.Total.Graded, r.Total.Criteria, r.Total.Strategy)
}

// num выводит балл без лишних нулей: 8, 7.5.
func num(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }

// total округляет итог до сотых только при выводе.
func total(v float64) string { return strconv.FormatFloat(v, 'f', 2, 64) }
