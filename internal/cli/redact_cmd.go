package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/redactor"
)

type redactArgs struct {
	file     string
	mask     string
	extraKeys string
	format   string
}

func parseRedactArgs(args []string) (redactArgs, error) {
	fs := flag.NewFlagSet("redact", flag.ContinueOnError)
	mask := fs.String("mask", "***", "replacement string for redacted values")
	extra := fs.String("keys", "", "comma-separated list of additional keys to redact")
	fmt_ := fs.String("format", "env", "output format: env or json")
	if err := fs.Parse(args); err != nil {
		return redactArgs{}, err
	}
	if fs.NArg() < 1 {
		return redactArgs{}, fmt.Errorf("usage: envdiff redact [flags] <file>")
	}
	return redactArgs{
		file:      fs.Arg(0),
		mask:      *mask,
		extraKeys: *extra,
		format:    *fmt_,
	}, nil
}

// RunRedact parses a .env file, redacts sensitive values, and writes the
// result to w.
func RunRedact(args []string, w io.Writer) error {
	a, err := parseRedactArgs(args)
	if err != nil {
		return err
	}

	env, err := parser.ParseFile(a.file)
	if err != nil {
		return fmt.Errorf("reading %s: %w", a.file, err)
	}

	opts := redactor.DefaultOptions()
	opts.Mask = a.mask
	if a.extraKeys != "" {
		for _, k := range strings.Split(a.extraKeys, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				opts.ExplicitKeys = append(opts.ExplicitKeys, k)
			}
		}
	}

	masked := redactor.Redact(env, opts)

	switch a.format {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(masked)
	default:
		for k, v := range masked {
			fmt.Fprintf(w, "%s=%s\n", k, v)
		}
		return nil
	}
}

// redactMain is the entry-point wired up in cmd/envdiff/main.go.
func redactMain(args []string) {
	if err := RunRedact(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
