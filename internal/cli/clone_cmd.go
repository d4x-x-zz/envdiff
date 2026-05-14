package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nicholasgasior/envdiff/internal/cloner"
	"github.com/nicholasgasior/envdiff/internal/parser"
)

type cloneArgs struct {
	input     string
	output    string
	prefix    string
	suffix    string
	uppercase bool
	omitEmpty bool
	format    string
}

func parseCloneArgs(args []string) (cloneArgs, error) {
	fs := flag.NewFlagSet("clone", flag.ContinueOnError)
	prefix := fs.String("prefix", "", "prefix to add to every key")
	suffix := fs.String("suffix", "", "suffix to add to every key")
	upper := fs.Bool("uppercase", false, "convert keys to uppercase")
	omit := fs.Bool("omit-empty", false, "omit keys with empty values")
	out := fs.String("out", "", "write output to file instead of stdout")
	fmt_ := fs.String("format", "dotenv", "output format: dotenv|json")
	if err := fs.Parse(args); err != nil {
		return cloneArgs{}, err
	}
	if fs.NArg() < 1 {
		return cloneArgs{}, fmt.Errorf("usage: envdiff clone [flags] <file>")
	}
	return cloneArgs{
		input:     fs.Arg(0),
		output:    *out,
		prefix:    *prefix,
		suffix:    *suffix,
		uppercase: *upper,
		omitEmpty: *omit,
		format:    *fmt_,
	}, nil
}

// RunClone parses a .env file, clones it with the requested transformations,
// and writes the result to stdout or a file.
func RunClone(args []string, stdout io.Writer) error {
	a, err := parseCloneArgs(args)
	if err != nil {
		return err
	}
	src, err := parser.ParseFile(a.input)
	if err != nil {
		return fmt.Errorf("parse %s: %w", a.input, err)
	}
	opts := cloner.Options{
		KeyPrefix:     a.prefix,
		KeySuffix:     a.suffix,
		UppercaseKeys: a.uppercase,
		OmitEmpty:     a.omitEmpty,
	}
	result := cloner.Clone(src, opts)

	var buf strings.Builder
	switch a.format {
	case "json":
		b, _ := json.MarshalIndent(result, "", "  ")
		buf.Write(b)
		buf.WriteByte('\n')
	default:
		for k, v := range result {
			fmt.Fprintf(&buf, "%s=%s\n", k, v)
		}
	}

	if a.output != "" {
		return os.WriteFile(a.output, []byte(buf.String()), 0o644)
	}
	_, err = fmt.Fprint(stdout, buf.String())
	return err
}
