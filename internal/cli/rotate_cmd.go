package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/rotator"
)

type rotateArgs struct {
	file     string
	renames  string // "OLD:NEW,OLD2:NEW2"
	strategy string
	output   string
}

func parseRotateArgs(args []string) (rotateArgs, error) {
	fs := flag.NewFlagSet("rotate", flag.ContinueOnError)
	file := fs.String("file", "", "path to .env file")
	renames := fs.String("renames", "", "comma-separated OLD:NEW pairs")
	strategy := fs.String("strategy", "remove", "old-key strategy: remove|deprecate|keep")
	output := fs.String("output", "", "write result to file (default: stdout)")
	if err := fs.Parse(args); err != nil {
		return rotateArgs{}, err
	}
	if *file == "" {
		return rotateArgs{}, fmt.Errorf("--file is required")
	}
	if *renames == "" {
		return rotateArgs{}, fmt.Errorf("--renames is required")
	}
	return rotateArgs{file: *file, renames: *renames, strategy: *strategy, output: *output}, nil
}

// RunRotate executes the rotate sub-command.
func RunRotate(args []string, stdout, _ *os.File) error {
	ra, err := parseRotateArgs(args)
	if err != nil {
		return err
	}

	env, err := parser.ParseFile(ra.file)
	if err != nil {
		return fmt.Errorf("rotate: %w", err)
	}

	var rotations []rotator.Rotation
	for _, pair := range strings.Split(ra.renames, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("rotate: invalid rename pair %q", pair)
		}
		rotations = append(rotations, rotator.Rotation{OldKey: parts[0], NewKey: parts[1]})
	}

	var strat rotator.Strategy
	switch ra.strategy {
	case "remove":
		strat = rotator.StrategyRemove
	case "deprecate":
		strat = rotator.StrategyDeprecate
	case "keep":
		strat = rotator.StrategyKeep
	default:
		return fmt.Errorf("rotate: unknown strategy %q", ra.strategy)
	}

	opts := rotator.Options{Strategy: strat, FailOnMissing: false}
	result, err := rotator.Rotate(env, rotations, opts)
	if err != nil {
		return err
	}

	var lines []string
	for k, v := range result {
		lines = append(lines, fmt.Sprintf("%s=%s", k, v))
	}
	output := strings.Join(lines, "\n") + "\n"

	if ra.output != "" {
		return os.WriteFile(ra.output, []byte(output), 0644)
	}
	fmt.Fprint(stdout, output)
	return nil
}
