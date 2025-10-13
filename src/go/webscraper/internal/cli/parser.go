package cli

import (
	"errors"
	"flag"
	"io"
	"strings"
)

const defaultConfigPath = "config/config.yaml"

// ErrHelp mirrors the standard flag.ErrHelp sentinel and allows callers to
// detect when users requested CLI usage output.
var ErrHelp = flag.ErrHelp

// Options represents command-line configuration passed to the scraper binary.
type Options struct {
	ConfigPath string
	URLs       []string
	Override   bool
}

// Parse converts raw CLI arguments into Options. It keeps positional arguments
// as URL overrides and applies sensible defaults for optional flags.
func Parse(args []string) (Options, error) {
	fs := flag.NewFlagSet("scraper", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	opts := Options{ConfigPath: defaultConfigPath}
	fs.StringVar(&opts.ConfigPath, "config", opts.ConfigPath, "Path to YAML config file")

	var directURLs multiValue
	fs.Var(&directURLs, "url", "Explicit URL to scrape (repeatable, overrides config input)")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return opts, ErrHelp
		}
		return opts, err
	}

	opts.URLs = append(opts.URLs, directURLs...)
	if len(directURLs) > 0 {
		opts.Override = true
	}
	opts.URLs = append(opts.URLs, fs.Args()...)
	opts.URLs = filterEmpty(opts.URLs)

	return opts, nil
}

// Usage returns a short CLI usage hint.
func Usage() string {
	return "scraper [--config path] [--url https://...] [url ...]"
}

type multiValue []string

func (m *multiValue) String() string {
	return strings.Join(*m, ",")
}

func (m *multiValue) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func filterEmpty(values []string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		if strings.TrimSpace(v) == "" {
			continue
		}
		out = append(out, v)
	}
	return out
}
