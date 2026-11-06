package cli

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type input struct {
	sc  *bufio.Scanner
	out io.Writer
}

func newInput(r io.Reader, w io.Writer) *input {
	return &input{bufio.NewScanner(r), w}
}

// line печатает подсказку и читает строку. Конец ввода (Ctrl+D) — io.EOF.
func (in *input) line(promt string) (string, error) {
	fmt.Fprintf(in.out, "%s: ", promt)
	if !in.sc.Scan() {
		if err := in.sc.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return strings.TrimSpace(in.sc.Text()), nil
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
