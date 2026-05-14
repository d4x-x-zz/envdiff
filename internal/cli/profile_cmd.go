package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/profiler"
)

type profileArgs struct {
	file   string
	format string
}

func parseProfileArgs(args []string) (profileArgs, error) {
	fs := flag.NewFlagSet("profile", flag.ContinueOnError)
	format := fs.String("format", "text", "output format: text or json")
	if err := fs.Parse(args); err != nil {
		return profileArgs{}, err
	}
	if fs.NArg() < 1 {
		return profileArgs{}, fmt.Errorf("usage: envdiff profile <file> [--format text|json]")
	}
	return profileArgs{file: fs.Arg(0), format: *format}, nil
}

// RunProfile parses a .env file and prints a statistical profile.
func RunProfile(args []string, out io.Writer) error {
	pa, err := parseProfileArgs(args)
	if err != nil {
		return err
	}
	env, err := parser.ParseFile(pa.file)
	if err != nil {
		return fmt.Errorf("reading %s: %w", pa.file, err)
	}
	p := profiler.Analyze(env, profiler.DefaultOptions())

	switch pa.format {
	case "json":
		return printProfileJSON(p, out)
	default:
		return printProfileText(p, out)
	}
}

func printProfileText(p profiler.Profile, out io.Writer) error {
	fmt.Fprintf(out, "Total keys   : %d\n", p.TotalKeys)
	fmt.Fprintf(out, "Empty values : %d\n", p.EmptyValues)
	fmt.Fprintf(out, "Density      : %.2f\n", p.Density)
	fmt.Fprintf(out, "Avg val len  : %.2f\n", p.AvgValueLen)
	fmt.Fprintf(out, "Max val len  : %d\n", p.MaxValueLen)
	fmt.Fprintf(out, "Min val len  : %d\n", p.MinValueLen)
	fmt.Fprintf(out, "Type breakdown:\n")
	for _, t := range profiler.SortedTypes(p) {
		fmt.Fprintf(out, "  %-8s: %d\n", t, p.TypeBreakdown[t])
	}
	return nil
}

func printProfileJSON(p profiler.Profile, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}

func profileMain() {
	if err := RunProfile(os.Args[2:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
