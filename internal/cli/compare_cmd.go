package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/parser"
)

type compareArgs struct {
	files  []string
	format string
}

func parseCompareArgs(args []string) (compareArgs, error) {
	fs := flag.NewFlagSet("compare", flag.ContinueOnError)
	format := fs.String("format", "text", "output format: text or json")
	if err := fs.Parse(args); err != nil {
		return compareArgs{}, err
	}
	if fs.NArg() < 2 {
		return compareArgs{}, fmt.Errorf("compare requires at least two files")
	}
	return compareArgs{files: fs.Args(), format: *format}, nil
}

// RunCompare executes the compare sub-command.
func RunCompare(args []string, out io.Writer) error {
	ca, err := parseCompareArgs(args)
	if err != nil {
		return err
	}

	fileEnv := comparator.FileEnv{}
	for _, path := range ca.files {
		env, err := parser.ParseFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}
		fileEnv[path] = env
	}

	report := comparator.Compare(fileEnv)

	switch strings.ToLower(ca.format) {
	case "json":
		return printCompareJSON(report, out)
	default:
		return printCompareText(report, out)
	}
}

func printCompareText(r *comparator.Report, out io.Writer) error {
	if r.TotalMissing() == 0 && r.TotalMismatched() == 0 {
		fmt.Fprintln(out, "✓ all keys match across files")
		return nil
	}
	for _, s := range r.Statuses {
		if len(s.Missing) > 0 {
			fmt.Fprintf(out, "MISSING  %s  (absent in: %s)\n", s.Key, strings.Join(s.Missing, ", "))
		} else if !s.Uniform {
			parts := make([]string, 0, len(s.Values))
			for f, v := range s.Values {
				parts = append(parts, fmt.Sprintf("%s=%s", f, v))
			}
			fmt.Fprintf(out, "MISMATCH %s  (%s)\n", s.Key, strings.Join(parts, " | "))
		}
	}
	return nil
}

func printCompareJSON(r *comparator.Report, out io.Writer) error {
	type entry struct {
		Key     string            `json:"key"`
		Status  string            `json:"status"`
		Values  map[string]string `json:"values,omitempty"`
		Missing []string          `json:"missing,omitempty"`
	}
	entries := []entry{}
	for _, s := range r.Statuses {
		if len(s.Missing) > 0 {
			entries = append(entries, entry{Key: s.Key, Status: "missing", Missing: s.Missing})
		} else if !s.Uniform {
			entries = append(entries, entry{Key: s.Key, Status: "mismatch", Values: s.Values})
		}
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

func compareMain(args []string) {
	if err := RunCompare(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
