package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/encoder"
	"github.com/user/envdiff/internal/parser"
)

type encodeArgs struct {
	file      string
	format    string
	omitEmpty bool
	output    string
}

func parseEncodeArgs(args []string) (encodeArgs, error) {
	fs := flag.NewFlagSet("encode", flag.ContinueOnError)
	format := fs.String("format", "shell", "output format: shell, docker, yaml")
	omitEmpty := fs.Bool("omit-empty", false, "omit keys with empty values")
	output := fs.String("o", "", "write output to file instead of stdout")

	if err := fs.Parse(args); err != nil {
		return encodeArgs{}, err
	}
	if fs.NArg() < 1 {
		return encodeArgs{}, fmt.Errorf("usage: envdiff encode [flags] <file>")
	}
	return encodeArgs{
		file:      fs.Arg(0),
		format:    *format,
		omitEmpty: *omitEmpty,
		output:    *output,
	}, nil
}

// RunEncode encodes a .env file into the requested serialization format.
func RunEncode(args []string, out io.Writer) error {
	a, err := parseEncodeArgs(args)
	if err != nil {
		return err
	}

	env, err := parser.ParseFile(a.file)
	if err != nil {
		return fmt.Errorf("parse %s: %w", a.file, err)
	}

	opts := encoder.DefaultOptions()
	opts.Format = encoder.Format(a.format)
	opts.OmitEmpty = a.omitEmpty

	result, err := encoder.Encode(env, opts)
	if err != nil {
		return err
	}

	if a.output != "" {
		if err := os.WriteFile(a.output, []byte(result), 0644); err != nil {
			return fmt.Errorf("write %s: %w", a.output, err)
		}
		fmt.Fprintf(out, "encoded output written to %s\n", a.output)
		return nil
	}

	fmt.Fprint(out, result)
	return nil
}
