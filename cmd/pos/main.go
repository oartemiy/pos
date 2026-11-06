package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"pos/internal/adapters/primary/cli"
	"pos/internal/application/facade"
	"pos/internal/di"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка:", err)
		os.Exit(1)
	}
}

func run() error {
	var cfg di.Config

	flag.StringVar(&cfg.DataDir, "data", "data", "каталог для хранения данных и выгрузок")
	flag.StringVar(&cfg.CourseTitle, "course", "КПО", "название курса (рубрикатора)")
	flag.StringVar(&cfg.Strategy, "strategy", "average", "правило итога: average | sum")
	flag.Parse()

	// Ctrl+C отменяет контекст, и операции ниже по стеку прерываются.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	container, err := di.Build(cfg)
	if err != nil {
		return err
	}

	return container.Invoke(func(f *facade.Facade) error {
		return cli.New(f, os.Stdin, os.Stdout).Run(ctx)
	})
}
