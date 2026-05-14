package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envdiff/internal/caster"
	"github.com/user/envdiff/internal/parser"
)

type castArgs struct {
	file   string
	format string
}

func parseCastArgs(args []string) (castArgs, error) {
	fs := flag.NewFlagSet("cast", flag.ContinueOnError)
	format := fs.String("format", "text", "output format: text or json")
	if err := fs.Parse(args); err != nil {
		return castArgs{}, err
	}
	if fs.NArg() < 1 {
		return castArgs{}, fmt.Errorf("usage: envdiff cast <file> [--format text|json]")
	}
	return castArgs{file: fs.Arg(0), format: *format}, nil
}

// RunCast parses a .env file and prints inferred types for each key.
func RunCast(args []string, out io.Writer) error {
	a, err := parseCastArgs(args)
	if err != nil {
		return err
	}

	env, err := parser.ParseFile(a.file)
	if err != nil {
		return fmt.Errorf("parse %s: %w", a.file, err)
	}

	results := caster.Cast(env, caster.DefaultOptions())

	// stable output
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	switch a.format {
	case "json":
		return renderCastJSON(results, out)
	default:
		return renderCastText(results, out)
	}
}

func renderCastText(results []caster.Result, out io.Writer) error {
	for _, r := range results {
		fmt.Fprintf(out, "%-30s %s\n", r.Key, r.Type)
	}
	return nil
}

func renderCastJSON(results []caster.Result, out io.Writer) error {
	type row struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Type  string `json:"type"`
	}
	rows := make([]row, len(results))
	for i, r := range results {
		rows[i] = row{Key: r.Key, Value: r.Value, Type: string(r.Type)}
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}

func init() {
	_ = os.Stderr // ensure os import used
}
