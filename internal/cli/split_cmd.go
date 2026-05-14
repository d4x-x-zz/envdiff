package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/your-org/envdiff/internal/parser"
	"github.com/your-org/envdiff/internal/splitter"
)

type splitArgs struct {
	file     string
	prefixes []string
	format   string
	keepFull bool
}

func parseSplitArgs(args []string) (splitArgs, error) {
	fs := flag.NewFlagSet("split", flag.ContinueOnError)
	prefixList := fs.String("prefixes", "", "comma-separated list of prefixes to split on")
	format := fs.String("format", "text", "output format: text or json")
	keepFull := fs.Bool("keep-full", false, "keep full key names (do not strip prefix)")

	if err := fs.Parse(args); err != nil {
		return splitArgs{}, err
	}
	if fs.NArg() < 1 {
		return splitArgs{}, fmt.Errorf("usage: envdiff split [flags] <file>")
	}

	var prefixes []string
	for _, p := range strings.Split(*prefixList, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			prefixes = append(prefixes, p)
		}
	}

	return splitArgs{
		file:     fs.Arg(0),
		prefixes: prefixes,
		format:   *format,
		keepFull: *keepFull,
	}, nil
}

// RunSplit executes the split sub-command.
func RunSplit(args []string, out io.Writer) error {
	a, err := parseSplitArgs(args)
	if err != nil {
		return err
	}

	env, err := parser.ParseFile(a.file)
	if err != nil {
		return fmt.Errorf("parse %s: %w", a.file, err)
	}

	opts := splitter.DefaultOptions()
	opts.Prefixes = a.prefixes
	opts.StripPrefix = !a.keepFull

	groups := splitter.Split(env, opts)

	if a.format == "json" {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(groups)
	}

	// text output
	for _, name := range splitter.SortedGroupNames(groups) {
		fmt.Fprintf(out, "[%s]\n", name)
		group := groups[name]
		keys := make([]string, 0, len(group))
		for k := range group {
			keys = append(keys, k)
		}
		for _, k := range keys {
			fmt.Fprintf(out, "  %s=%s\n", k, group[k])
		}
	}
	return nil
}

func splitMain() {
	if err := RunSplit(os.Args[2:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
