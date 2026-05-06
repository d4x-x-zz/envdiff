package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/reporter"
)

// Config holds parsed CLI flags and arguments.
type Config struct {
	Format string
	Strict bool
	LeftFile  string
	RightFile string
}

// Run parses args and executes the diff workflow.
func Run(args []string) error {
	cfg, err := parseArgs(args)
	if err != nil {
		return err
	}

	left, err := parser.ParseFile(cfg.LeftFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", cfg.LeftFile, err)
	}

	right, err := parser.ParseFile(cfg.RightFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", cfg.RightFile, err)
	}

	result := differ.Diff(left, right)

	r := reporter.New(os.Stdout, cfg.Format)
	if err := r.Report(result, cfg.LeftFile, cfg.RightFile); err != nil {
		return fmt.Errorf("reporting: %w", err)
	}

	if cfg.Strict && !result.Clean() {
		return errors.New("differences found")
	}

	return nil
}

func parseArgs(args []string) (*Config, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)
	format := fs.String("format", "text", "output format: text or json")
	strict := fs.Bool("strict", false, "exit with non-zero status if differences found")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if fs.NArg() != 2 {
		return nil, fmt.Errorf("usage: envdiff [flags] <file1> <file2>")
	}

	if *format != "text" && *format != "json" {
		return nil, fmt.Errorf("unknown format %q: must be text or json", *format)
	}

	return &Config{
		Format:    *format,
		Strict:    *strict,
		LeftFile:  fs.Arg(0),
		RightFile: fs.Arg(1),
	}, nil
}
