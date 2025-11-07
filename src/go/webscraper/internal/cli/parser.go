package cli

import (
	"errors"
	"flag"
	"io"
)

const defaultConfigPath = "config/config.yaml"

// ErrHelp mirrors the standard flag.ErrHelp sentinel and allows callers to
// detect when users requested CLI usage output.
var ErrHelp = flag.ErrHelp

// Options represents command-line configuration passed to the scraper binary.
type Options struct {
	ConfigPath string
}

// Parse converts raw CLI arguments into Options. It keeps positional arguments
// as URL overrides and applies sensible defaults for optional flags.
func Parse(args []string) (Options, error) {
	fs := flag.NewFlagSet("scraper", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	opts := Options{ConfigPath: defaultConfigPath}
	fs.StringVar(&opts.ConfigPath, "config", opts.ConfigPath, "Path to YAML config file")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return opts, ErrHelp
		}
		return opts, err
	}

	return opts, nil
}

// Usage returns a short CLI usage hint.
func Usage() string {
	return "scraper [--config path]"
}
