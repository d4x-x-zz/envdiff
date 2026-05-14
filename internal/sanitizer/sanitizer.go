// Package sanitizer removes or replaces characters in env values
// that are unsafe for specific target environments (e.g. shell, docker).
package sanitizer

import (
	"strings"
	"unicode"
)

// Options controls sanitizer behaviour.
type Options struct {
	// StripControlChars removes non-printable / control characters from values.
	StripControlChars bool
	// TrimWhitespace trims leading and trailing whitespace from values.
	TrimWhitespace bool
	// ReplaceNewlines replaces embedded newline characters with the given string.
	// If empty, newlines are removed.
	ReplaceNewlines *string
	// MaxLength truncates values to this length. 0 means no limit.
	MaxLength int
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	nl := `\n`
	return Options{
		StripControlChars: true,
		TrimWhitespace:    true,
		ReplaceNewlines:   &nl,
		MaxLength:         0,
	}
}

// Sanitize applies the given options to every value in env and returns a new map.
func Sanitize(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = sanitizeValue(v, opts)
	}
	return out
}

func sanitizeValue(v string, opts Options) string {
	if opts.TrimWhitespace {
		v = strings.TrimSpace(v)
	}

	if opts.ReplaceNewlines != nil {
		v = strings.ReplaceAll(v, "\n", *opts.ReplaceNewlines)
		v = strings.ReplaceAll(v, "\r", "")
	}

	if opts.StripControlChars {
		var b strings.Builder
		for _, r := range v {
			if r == '\t' || !unicode.IsControl(r) {
				b.WriteRune(r)
			}
		}
		v = b.String()
	}

	if opts.MaxLength > 0 && len(v) > opts.MaxLength {
		v = v[:opts.MaxLength]
	}

	return v
}
