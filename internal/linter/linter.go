// Package linter provides checks for common .env file style issues.
package linter

import (
	"fmt"
	"strings"
	"unicode"
)

// Issue represents a single linting problem found in an env map.
type Issue struct {
	Key     string
	Message string
}

// Options controls which lint checks are enabled.
type Options struct {
	CheckUpperCase    bool // keys should be UPPER_CASE
	CheckNoSpaces     bool // keys should not contain spaces
	CheckNoLeadDigit  bool // keys should not start with a digit
	CheckNoDupPrefix  bool // warn on duplicate key prefixes (e.g. APP_ appears 10+ times)
}

// DefaultOptions returns a sensible default lint configuration.
func DefaultOptions() Options {
	return Options{
		CheckUpperCase:   true,
		CheckNoSpaces:    true,
		CheckNoLeadDigit: true,
		CheckNoDupPrefix: false,
	}
}

// Lint runs all enabled checks against the provided env map and returns
// a slice of issues. An empty slice means no problems were found.
func Lint(env map[string]string, opts Options) []Issue {
	var issues []Issue

	for key := range env {
		if opts.CheckNoSpaces && strings.ContainsAny(key, " \t") {
			issues = append(issues, Issue{Key: key, Message: "key contains whitespace"})
		}

		if opts.CheckNoLeadDigit && len(key) > 0 && unicode.IsDigit(rune(key[0])) {
			issues = append(issues, Issue{Key: key, Message: "key starts with a digit"})
		}

		if opts.CheckUpperCase && key != strings.ToUpper(key) {
			issues = append(issues, Issue{
				Key:     key,
				Message: fmt.Sprintf("key is not upper-case (got %q, want %q)", key, strings.ToUpper(key)),
			})
		}
	}

	if opts.CheckNoDupPrefix {
		issues = append(issues, checkDupPrefixes(env)...)
	}

	return issues
}

// checkDupPrefixes warns when a single PREFIX_ accounts for more than half the keys.
func checkDupPrefixes(env map[string]string) []Issue {
	prefixCount := map[string]int{}
	for key := range env {
		parts := strings.SplitN(key, "_", 2)
		if len(parts) == 2 {
			prefixCount[parts[0]]++
		}
	}
	total := len(env)
	var issues []Issue
	for prefix, count := range prefixCount {
		if total > 0 && count*2 > total {
			issues = append(issues, Issue{
				Key:     prefix + "_*",
				Message: fmt.Sprintf("prefix %q dominates the file (%d/%d keys)", prefix, count, total),
			})
		}
	}
	return issues
}
