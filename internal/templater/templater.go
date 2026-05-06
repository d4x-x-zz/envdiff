// Package templater generates a .env.template file from an existing env map,
// replacing all values with empty strings or typed placeholders.
package templater

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls how the template is generated.
type Options struct {
	// UseTypedPlaceholders replaces values with type hints like <string>, <number>.
	UseTypedPlaceholders bool
	// CommentOriginal includes the original value as a comment above each key.
	CommentOriginal bool
}

// DefaultOptions returns sensible defaults for template generation.
func DefaultOptions() Options {
	return Options{
		UseTypedPlaceholders: true,
		CommentOriginal:      false,
	}
}

// Generate takes an env map and returns a template string.
func Generate(env map[string]string, opts Options) string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := env[k]
		if opts.CommentOriginal {
			fmt.Fprintf(&sb, "# original: %s\n", v)
		}
		placeholder := ""
		if opts.UseTypedPlaceholders {
			placeholder = inferPlaceholder(v)
		}
		fmt.Fprintf(&sb, "%s=%s\n", k, placeholder)
	}
	return sb.String()
}

// inferPlaceholder returns a type hint based on the original value.
func inferPlaceholder(v string) string {
	if v == "" {
		return "<string>"
	}
	if isNumeric(v) {
		return "<number>"
	}
	if v == "true" || v == "false" {
		return "<bool>"
	}
	return "<string>"
}

// isNumeric returns true if the string looks like an integer or float.
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	dot := false
	for i, c := range s {
		if c == '-' && i == 0 {
			continue
		}
		if c == '.' && !dot {
			dot = true
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
