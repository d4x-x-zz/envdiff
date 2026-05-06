package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/linter"
	"github.com/user/envdiff/internal/parser"
)

type lintArgs struct {
	file          string
	allowLower    bool
	allowSpaces   bool
	allowDigitStart bool
	allowDupPrefix bool
}

func parseLintArgs(args []string) (lintArgs, error) {
	fs := flag.NewFlagSet("lint", flag.ContinueOnError)

	allowLower := fs.Bool("allow-lower", false, "allow lowercase key names")
	allowSpaces := fs.Bool("allow-spaces", false, "allow spaces in key names")
	allowDigitStart := fs.Bool("allow-digit-start", false, "allow keys starting with a digit")
	allowDupPrefix := fs.Bool("allow-dup-prefix", false, "allow duplicate key prefixes")

	if err := fs.Parse(args); err != nil {
		return lintArgs{}, err
	}

	if fs.NArg() < 1 {
		return lintArgs{}, fmt.Errorf("usage: envdiff lint [options] <file>")
	}

	return lintArgs{
		file:            fs.Arg(0),
		allowLower:      *allowLower,
		allowSpaces:     *allowSpaces,
		allowDigitStart: *allowDigitStart,
		allowDupPrefix:  *allowDupPrefix,
	}, nil
}

func RunLint(args []string, out *os.File) error {
	la, err := parseLintArgs(args)
	if err != nil {
		return err
	}

	env, err := parser.ParseFile(la.file)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", la.file, err)
	}

	opts := linter.DefaultOptions()
	if la.allowLower {
		opts.RequireUpperCase = false
	}
	if la.allowSpaces {
		opts.DisallowSpaces = false
	}
	if la.allowDigitStart {
		opts.DisallowDigitStart = false
	}
	if la.allowDupPrefix {
		opts.CheckDupPrefixes = false
	}

	issues := linter.Lint(env, opts)
	if len(issues) == 0 {
		fmt.Fprintln(out, "no lint issues found")
		return nil
	}

	for _, issue := range issues {
		fmt.Fprintf(out, "[%s] %s\n", issue.Key, issue.Message)
	}

	return fmt.Errorf("%d lint issue(s) found", len(issues))
}
