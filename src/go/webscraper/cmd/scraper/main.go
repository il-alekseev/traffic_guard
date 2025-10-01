package main

import (
	"context"
	"errors"
	"fmt"
	"os"

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

	if options.Override {
		cfg.InputPath = ""
	}

	params := bootstrap.Params{
		Config: cfg,
		URLs:   options.URLs,
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
	}

	if err := bootstrap.Run(context.Background(), params); err != nil {
		if errors.Is(err, bootstrap.ErrNoURLs) {
			return fmt.Errorf("no urls provided: configure input_path or pass urls")
		}
		return err
	}

	return nil
}
