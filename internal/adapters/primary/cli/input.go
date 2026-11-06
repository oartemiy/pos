package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type input struct {
	r     io.Reader
	out   io.Writer
	lines chan readResult
	done  <-chan struct{} // закрывается по Ctrl+C; nil — ожидание без отмены
}

type readResult struct {
	text string
	err  error
}

func newInput(r io.Reader, w io.Writer) *input {
	return &input{r: r, out: w}
}

// bind связывает ввод с контекстом: после его отмены line сразу
// возвращает ошибку контекста, не дожидаясь, пока пользователь нажмёт Enter.
func (in *input) bind(ctx context.Context) {
	in.done = ctx.Done()
}

// start запускает чтение в отдельной горутине: блокирующий Scan
// нельзя прервать, а чтение из канала — можно.
func (in *input) start() {
	in.lines = make(chan readResult)
	go func() {
		defer close(in.lines)
		sc := bufio.NewScanner(in.r)
		for sc.Scan() {
			in.lines <- readResult{text: sc.Text()}
		}
		err := sc.Err()
		if err == nil {
			err = io.EOF
		}
		in.lines <- readResult{err: err}
	}()
}

// line печатает подсказку и читает строку. Конец ввода (Ctrl+D) — io.EOF,
// отмена контекста (Ctrl+C) — context.Canceled.
func (in *input) line(promt string) (string, error) {
	if in.lines == nil {
		in.start()
	}
	fmt.Fprintf(in.out, "%s: ", promt)
	select {
	case res, ok := <-in.lines:
		if !ok {
			return "", io.EOF
		}
		if res.err != nil {
			return "", res.err
		}
		return strings.TrimSpace(res.text), nil
	case <-in.done:
		return "", context.Canceled
	}
}

// number читает число, принимая и точку, и запятую. Пустой ввод — отмена.
func (in *input) number(promt string) (float64, error) {
	for {
		s, err := in.line(promt)
		if err != nil {
			return 0, err
		}
		if s == "" {
			return 0, errCancelled
		}
		v, err := strconv.ParseFloat(strings.Replace(s, ",", ".", 1), 64)
		if err == nil {
			return v, nil
		}
		fmt.Fprintln(in.out, "  Введите число, например 7; 7,5 или 7.5")
	}
}

// index читает номер строки из списка (1..n) и возвращает индекс с нуля.
func (in *input) index(prompt string, n int) (int, error) {
	for {
		s, err := in.line(fmt.Sprintf("%s (1–%d, Enter — отмена)", prompt, n))
		if err != nil {
			return 0, err
		}
		if s == "" {
			return 0, errCancelled
		}
		i, err := strconv.Atoi(s)
		if err == nil && i >= 1 && i <= n {
			return i - 1, nil
		}
		fmt.Fprintf(in.out, "  Введите номер от 1 до %d\n", n)
	}
}

func (in *input) confirm(prompt string) (bool, error) {
	s, err := in.line(prompt + " (да/нет)")
	if err != nil {
		return false, err
	}
	switch strings.ToLower(s) {
	case "да", "д", "yes", "y":
		return true, nil
	default:
		return false, nil
	}
}
