package cli

import (
	"fmt"
	"os"
	"strings"

	"envdiff/internal/merger"
	"envdiff/internal/parser"
)

// MergeArgs holds parsed arguments for the merge sub-command.
type MergeArgs struct {
	Files    []string
	Strategy merger.Strategy
	Output   string // "stdout" or a file path
}

// parseMergeArgs parses os.Args for the merge sub-command.
// Expected form: envdiff merge [--strategy=first|last|error] [--out=FILE] file1 file2 ...
func parseMergeArgs(args []string) (MergeArgs, error) {
	ma := MergeArgs{Strategy: merger.StrategyFirst, Output: "stdout"}

	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "--strategy="):
			val := strings.TrimPrefix(arg, "--strategy=")
			switch val {
			case "first":
				ma.Strategy = merger.StrategyFirst
			case "last":
				ma.Strategy = merger.StrategyLast
			case "error":
				ma.Strategy = merger.StrategyError
			default:
				return ma, fmt.Errorf("unknown strategy %q", val)
			}
		case strings.HasPrefix(arg, "--out="):
			ma.Output = strings.TrimPrefix(arg, "--out=")
		default:
			ma.Files = append(ma.Files, arg)
		}
	}

	if len(ma.Files) < 2 {
		return ma, fmt.Errorf("merge requires at least two files")
	}
	return ma, nil
}

// RunMerge executes the merge sub-command.
func RunMerge(args []string) error {
	ma, err := parseMergeArgs(args)
	if err != nil {
		return err
	}

	var maps []map[string]string
	for _, f := range ma.Files {
		m, err := parser.ParseFile(f)
		if err != nil {
			return fmt.Errorf("reading %s: %w", f, err)
		}
		maps = append(maps, m)
	}

	result, err := merger.Merge(maps, merger.Options{Strategy: ma.Strategy})
	if err != nil {
		return err
	}

	var lines []string
	for k, v := range result {
		lines = append(lines, fmt.Sprintf("%s=%s", k, v))
	}
	output := strings.Join(lines, "\n") + "\n"

	if ma.Output == "stdout" {
		fmt.Print(output)
		return nil
	}
	return os.WriteFile(ma.Output, []byte(output), 0o644)
}
