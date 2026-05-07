package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/prompts"
)

type promptArgs struct {
	file       string
	output     string
	skipFilled bool
	format     string
}

func parsePromptArgs(args []string) (promptArgs, error) {
	fs := flag.NewFlagSet("prompt", flag.ContinueOnError)
	output := fs.String("output", "", "write result to file instead of stdout")
	skip := fs.Bool("skip-filled", true, "skip keys that already have non-placeholder values")
	fmt_ := fs.String("format", "dotenv", "output format: dotenv|json|markdown")
	if err := fs.Parse(args); err != nil {
		return promptArgs{}, err
	}
	if fs.NArg() < 1 {
		return promptArgs{}, fmt.Errorf("usage: envdiff prompt [flags] <file>")
	}
	return promptArgs{
		file:       fs.Arg(0),
		output:     *output,
		skipFilled: *skip,
		format:     *fmt_,
	}, nil
}

// RunPrompt interactively fills missing/placeholder values in a .env file and
// writes the completed result to stdout or a file.
func RunPrompt(args []string, stdout, stderr *os.File) int {
	pa, err := parsePromptArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	env, err := parser.ParseFile(pa.file)
	if err != nil {
		fmt.Fprintf(stderr, "error reading %s: %v\n", pa.file, err)
		return 1
	}

	opts := prompts.DefaultOptions()
	opts.In = os.Stdin
	opts.Out = stderr // prompts go to stderr so stdout stays clean
	opts.SkipFilled = pa.skipFilled

	filled, err := prompts.Fill(env, opts)
	if err != nil {
		fmt.Fprintf(stderr, "prompt error: %v\n", err)
		return 1
	}

	expOpts := exporter.DefaultOptions()
	expOpts.Format = pa.format
	expOpts.Redact = false

	out, err := exporter.Export(filled, nil, expOpts)
	if err != nil {
		fmt.Fprintf(stderr, "export error: %v\n", err)
		return 1
	}

	if pa.output != "" {
		if err := os.WriteFile(pa.output, []byte(out), 0o644); err != nil {
			fmt.Fprintf(stderr, "write error: %v\n", err)
			return 1
		}
		fmt.Fprintf(stderr, "written to %s\n", pa.output)
	} else {
		fmt.Fprint(stdout, out)
	}
	return 0
}
