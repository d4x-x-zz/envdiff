// Package validator checks env maps for common issues such as empty values,
// keys that look like placeholders, and keys violating naming conventions.
package validator

import (
	"fmt"
	"strings"
)

// Issue represents a single validation warning for a key.
type Issue struct {
	Key     string
	Message string
}

func (i Issue) String() string {
	return fmt.Sprintf("%s: %s", i.Key, i.Message)
}

// Options controls which checks are enabled.
type Options struct {
	WarnEmpty       bool // flag keys with empty values
	WarnPlaceholder bool // flag values that look like <VALUE> or CHANGEME
	WarnNaming      bool // flag keys that contain lowercase letters
}

// DefaultOptions returns a sensible default set of checks.
func DefaultOptions() Options {
	return Options{
		WarnEmpty:       true,
		WarnPlaceholder: true,
		WarnNaming:      false,
	}
}

var placeholderHints = []string{
	"changeme",
	"todo",
	"fixme",
	"<value>",
	"your_",
	"example",
}

// Validate runs the enabled checks against the provided env map and returns
// any issues found. The map is keyed by env var name.
func Validate(env map[string]string, opts Options) []Issue {
	var issues []Issue

	for k, v := range env {
		if opts.WarnEmpty && strings.TrimSpace(v) == "" {
			issues = append(issues, Issue{Key: k, Message: "value is empty"})
			continue
		}

		if opts.WarnPlaceholder {
			lower := strings.ToLower(v)
			for _, hint := range placeholderHints {
				if strings.Contains(lower, hint) {
					issues = append(issues, Issue{
						Key:     k,
						Message: fmt.Sprintf("value looks like a placeholder (%q)", v),
					})
					break
				}
			}
		}

		if opts.WarnNaming && k != strings.ToUpper(k) {
			issues = append(issues, Issue{
				Key:     k,
				Message: "key contains lowercase letters",
			})
		}
	}

	return issues
}
