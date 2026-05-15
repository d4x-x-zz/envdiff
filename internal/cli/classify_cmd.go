package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envdiff/internal/classifier"
	"github.com/user/envdiff/internal/parser"
)

func parseClassifyArgs(args []string) (file string, jsonOut bool, err error) {
	fs := flag.NewFlagSet("classify", flag.ContinueOnError)
	fs.Bool("include-other", true, "include unmatched keys in 'other' category")
	var jFlag bool
	fs.BoolVar(&jFlag, "json", false, "output as JSON")
	if err = fs.Parse(args); err != nil {
		return
	}
	if fs.NArg() < 1 {
		err = fmt.Errorf("usage: envdiff classify [--json] <file>")
		return
	}
	return fs.Arg(0), jFlag, nil
}

// RunClassify implements the classify sub-command.
func RunClassify(args []string, out io.Writer) error {
	file, jsonOut, err := parseClassifyArgs(args)
	if err != nil {
		return err
	}

	env, err := parser.ParseFile(file)
	if err != nil {
		return fmt.Errorf("parse %s: %w", file, err)
	}

	opts := classifier.DefaultOptions()
	res := classifier.Classify(env, opts)

	if jsonOut {
		// convert map keys to strings for JSON
		plain := make(map[string][]string)
		for cat, keys := range res.Categories {
			sorted := make([]string, len(keys))
			copy(sorted, keys)
			sort.Strings(sorted)
			plain[string(cat)] = sorted
		}
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(plain)
	}

	cats := make([]string, 0, len(res.Categories))
	for cat := range res.Categories {
		cats = append(cats, string(cat))
	}
	sort.Strings(cats)

	for _, cat := range cats {
		keys := res.Categories[classifier.Category(cat)]
		sort.Strings(keys)
		fmt.Fprintf(out, "[%s]\n", cat)
		for _, k := range keys {
			fmt.Fprintf(out, "  %s\n", k)
		}
	}
	return nil
}

func classifyMain() {
	if err := RunClassify(os.Args[2:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
