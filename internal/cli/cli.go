package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/envdiff/internal/differ"
	"github.com/your-org/envdiff/internal/filter"
	"github.com/your-org/envdiff/internal/parser"
	"github.com/your-org/envdiff/internal/reporter"
)

// args holds parsed CLI arguments.
type args struct {
	leftFile       string
	rightFile      string
	strict         bool
	format         string
	onlyMissing    bool
	onlyMismatched bool
	keyPrefix      string
}

// Run is the entry point for the CLI.
func Run(argv []string) int {
	a, err := parseArgs(argv)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}

	left, err := parser.ParseFile(a.leftFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read %s: %v\n", a.leftFile, err)
		return 2
	}

	right, err := parser.ParseFile(a.rightFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read %s: %v\n", a.rightFile, err)
		return 2
	}

	result := differ.Diff(left, right)

	result = filter.Apply(result, filter.Options{
		OnlyMissing:    a.onlyMissing,
		OnlyMismatched: a.onlyMismatched,
		KeyPrefix:      a.keyPrefix,
	})

	r := reporter.New(os.Stdout, a.format)
	if err := r.Report(result); err != nil {
		fmt.Fprintln(os.Stderr, "report error:", err)
		return 2
	}

	if a.strict && !result.Clean() {
		return 1
	}
	return 0
}

func parseArgs(argv []string) (args, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)
	strict := fs.Bool("strict", false, "exit 1 when differences are found")
	format := fs.String("format", "text", "output format: text or json")
	onlyMissing := fs.Bool("only-missing", false, "show only missing keys")
	onlyMismatched := fs.Bool("only-mismatched", false, "show only mismatched values")
	keyPrefix := fs.String("prefix", "", "filter keys by prefix")

	if err := fs.Parse(argv); err != nil {
		return args{}, err
	}
	if fs.NArg() < 2 {
		return args{}, fmt.Errorf("usage: envdiff [flags] <left.env> <right.env>")
	}
	return args{
		leftFile:       fs.Arg(0),
		rightFile:      fs.Arg(1),
		strict:         *strict,
		format:         *format,
		onlyMissing:    *onlyMissing,
		onlyMismatched: *onlyMismatched,
		keyPrefix:      *keyPrefix,
	}, nil
}
