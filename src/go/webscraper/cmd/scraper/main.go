package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"scrapper/config"
	"scrapper/internal/bootstrap"
	"scrapper/internal/cli"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	options, err := cli.Parse(os.Args[1:])
	if err != nil {
		if errors.Is(err, cli.ErrHelp) {
			fmt.Println(cli.Usage())
			return nil
		}
		return fmt.Errorf("parse cli: %w", err)
	}

	cfg, err := config.Load(options.ConfigPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	config.ApplyEnvOverrides(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	params := bootstrap.Params{
		Config: cfg,
	}

	if err := bootstrap.Run(ctx, params); err != nil {
		return err
	}

	return nil
}
