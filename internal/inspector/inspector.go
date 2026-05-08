// Package inspector provides utilities for inspecting .env files and
// producing a summary of key statistics such as total keys, empty values,
// sensitive keys, and placeholder counts.
package inspector

import "strings"

// Summary holds the inspection results for a single env map.
type Summary struct {
	TotalKeys    int
	EmptyValues  int
	SensitiveKeys []string
	Placeholders int
	UniqueValues int
}

// Options controls how inspection is performed.
type Options struct {
	// SensitivePatterns are substrings that mark a key as sensitive.
	SensitivePatterns []string
	// PlaceholderPrefixes are value prefixes treated as placeholders (e.g. "<", "CHANGE_ME").
	PlaceholderPrefixes []string
}

// DefaultOptions returns sensible defaults for inspection.
func DefaultOptions() Options {
	return Options{
		SensitivePatterns:   []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PRIVATE"},
		PlaceholderPrefixes: []string{"<", "CHANGE_ME", "YOUR_", "TODO"},
	}
}

// Inspect analyses env and returns a Summary.
func Inspect(env map[string]string, opts Options) Summary {
	seen := make(map[string]struct{})
	s := Summary{TotalKeys: len(env)}

	for k, v := range env {
		if v == "" {
			s.EmptyValues++
		}
		if isSensitive(k, opts.SensitivePatterns) {
			s.SensitiveKeys = append(s.SensitiveKeys, k)
		}
		if isPlaceholder(v, opts.PlaceholderPrefixes) {
			s.Placeholders++
		}
		if v != "" {
			seen[v] = struct{}{}
		}
	}

	s.UniqueValues = len(seen)
	return s
}

func isSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

func isPlaceholder(value string, prefixes []string) bool {
	upper := strings.ToUpper(value)
	for _, p := range prefixes {
		if strings.HasPrefix(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}
