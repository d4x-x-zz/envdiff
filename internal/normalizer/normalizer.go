// Package normalizer standardises .env key and value formatting.
// It can uppercase keys, trim surrounding whitespace from values,
// and strip redundant surrounding quotes that some tools leave behind.
package normalizer

import (
	"strings"
)

// Options controls which normalisation passes are applied.
type Options struct {
	// UppercaseKeys converts every key to UPPER_CASE.
	UppercaseKeys bool
	// TrimValues strips leading/trailing whitespace from values.
	TrimValues bool
	// StripQuotes removes a single layer of matching surrounding quotes.
	StripQuotes bool
}

// DefaultOptions returns a sensible out-of-the-box configuration.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys: true,
		TrimValues:    true,
		StripQuotes:   false,
	}
}

// Normalize applies the configured passes to a copy of env and returns it.
// The original map is never mutated.
func Normalize(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		nk := k
		nv := v

		if opts.UppercaseKeys {
			nk = strings.ToUpper(nk)
		}
		if opts.TrimValues {
			nv = strings.TrimSpace(nv)
		}
		if opts.StripQuotes {
			nv = stripQuotes(nv)
		}

		out[nk] = nv
	}
	return out
}

// stripQuotes removes one layer of matching " or ' from both ends of s.
func stripQuotes(s string) string {
	if len(s) < 2 {
		return s
	}
	if (s[0] == '"' && s[len(s)-1] == '"') ||
		(s[0] == '\'' && s[len(s)-1] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}
