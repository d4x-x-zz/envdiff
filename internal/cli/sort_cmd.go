package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"envdiff/internal/parser"
	"envdiff/internal/sorter"
)

type sortArgs struct {
	file        string
	groupPrefix bool
	format      string // "text" or "json"
}

func parseSortArgs(args []string) (sortArgs, error) {
	sa := sortArgs{format: "text"}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--group-prefix":
			sa.groupPrefix = true
		case "--format":
			if i+1 >= len(args) {
				return sa, fmt.Errorf("--format requires a value")
			}
			i++
			sa.format = args[i]
		default:
			if sa.file == "" {
				sa.file = args[i]
			}
		}
	}
	if sa.file == "" {
		return sa, fmt.Errorf("a .env file path is required")
	}
	return sa, nil
}

// RunSort parses a .env file and prints its keys in sorted order.
func RunSort(args []string, out *os.File) error {
	sa, err := parseSortArgs(args)
	if err != nil {
		return err
	}

	env, err := parser.ParseFile(sa.file)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", sa.file, err)
	}

	opts := sorter.Options{
		Alphabetical:  true,
		GroupByPrefix: sa.groupPrefix,
	}
	keys := sorter.Sort(env, opts)

	switch strings.ToLower(sa.format) {
	case "json":
		ordered := make([]map[string]string, 0, len(keys))
		for _, k := range keys {
			ordered = append(ordered, map[string]string{"key": k, "value": env[k]})
		}
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(ordered)
	default:
		for _, k := range keys {
			fmt.Fprintf(out, "%s=%s\n", k, env[k])
		}
	}
	return nil
}
