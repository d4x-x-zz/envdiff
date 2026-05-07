// Package scoper provides functionality to scope (filter/extract) env maps
// to a specific namespace prefix, stripping or preserving the prefix.
package scoper

import (
	"sort"
	"strings"
)

// Options controls Scope behaviour.
type Options struct {
	// Prefix is the namespace prefix to scope to (e.g. "APP_").
	Prefix string
	// StripPrefix removes the prefix from keys in the output when true.
	StripPrefix bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		StripPrefix: true,
	}
}

// ScopedResult holds the output of a Scope call.
type ScopedResult struct {
	// Scoped contains the keys that matched the prefix.
	Scoped map[string]string
	// Excluded contains the keys that did not match the prefix.
	Excluded map[string]string
}

// Scope partitions env into keys that match opts.Prefix and those that don't.
// If opts.StripPrefix is true, the prefix is removed from matched keys.
func Scope(env map[string]string, opts Options) ScopedResult {
	scoped := make(map[string]string)
	excluded := make(map[string]string)

	for k, v := range env {
		if opts.Prefix == "" || strings.HasPrefix(k, opts.Prefix) {
			outKey := k
			if opts.StripPrefix && opts.Prefix != "" {
				outKey = strings.TrimPrefix(k, opts.Prefix)
			}
			scoped[outKey] = v
		} else {
			excluded[k] = v
		}
	}

	return ScopedResult{Scoped: scoped, Excluded: excluded}
}

// SortedKeys returns the keys of m in sorted order.
func SortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
