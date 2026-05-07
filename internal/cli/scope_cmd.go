package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/scoper"
)

type scopeArgs struct {
	file        string
	prefix      string
	stripPrefix bool
	format      string
	output      string
}

func parseScopeArgs(args []string) (scopeArgs, error) {
	fs := flag.NewFlagSet("scope", flag.ContinueOnError)
	prefix := fs.String("prefix", "", "namespace prefix to scope to")
	strip := fs.Bool("strip", true, "strip prefix from output keys")
	format := fs.String("format", "env", "output format: env or json")
	output := fs.String("out", "", "write output to file instead of stdout")

	if err := fs.Parse(args); err != nil {
		return scopeArgs{}, err
	}
	if fs.NArg() < 1 {
		return scopeArgs{}, fmt.Errorf("usage: envdiff scope -prefix PREFIX [flags] <file>")
	}
	return scopeArgs{
		file:        fs.Arg(0),
		prefix:      *prefix,
		stripPrefix: *strip,
		format:      *format,
		output:      *output,
	}, nil
}

// RunScope executes the scope sub-command.
func RunScope(args []string, stdout io.Writer) error {
	a, err := parseScopeArgs(args)
	if err != nil {
		return err
	}

	env, err := parser.ParseFile(a.file)
	if err != nil {
		return fmt.Errorf("parse %s: %w", a.file, err)
	}

	opts := scoper.Options{Prefix: a.prefix, StripPrefix: a.stripPrefix}
	res := scoper.Scope(env, opts)

	var sb strings.Builder
	switch a.format {
	case "json":
		b, jerr := json.MarshalIndent(map[string]any{
			"scoped":   res.Scoped,
			"excluded": res.Excluded,
		}, "", "  ")
		if jerr != nil {
			return jerr
		}
		sb.Write(b)
		sb.WriteByte('\n')
	default:
		for _, k := range scoper.SortedKeys(res.Scoped) {
			fmt.Fprintf(&sb, "%s=%s\n", k, res.Scoped[k])
		}
	}

	if a.output != "" {
		if werr := os.WriteFile(a.output, []byte(sb.String()), 0o644); werr != nil {
			return fmt.Errorf("write %s: %w", a.output, werr)
		}
		fmt.Fprintf(stdout, "scoped output written to %s\n", a.output)
		return nil
	}

	fmt.Fprint(stdout, sb.String())
	return nil
}
