// Package sorter provides utilities for sorting and grouping env file keys.
package sorter

import (
	"sort"
	"strings"
)

// Options controls how keys are sorted.
type Options struct {
	// Alphabetical sorts keys A-Z when true.
	Alphabetical bool
	// GroupByPrefix groups keys sharing the same prefix (e.g. DB_, AWS_).
	GroupByPrefix bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Alphabetical:  true,
		GroupByPrefix: false,
	}
}

// Sort returns a sorted slice of keys from the provided map according to opts.
func Sort(env map[string]string, opts Options) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	if opts.GroupByPrefix {
		return groupByPrefix(keys)
	}

	if opts.Alphabetical {
		sort.Strings(keys)
	}

	return keys
}

// groupByPrefix sorts keys so that keys sharing a common prefix
// (the part before the first underscore) appear together, and within
// each group keys are sorted alphabetically.
func groupByPrefix(keys []string) []string {
	groups := make(map[string][]string)
	order := []string{}
	seen := make(map[string]bool)

	for _, k := range keys {
		prefix := extractPrefix(k)
		if !seen[prefix] {
			seen[prefix] = true
			order = append(order, prefix)
		}
		groups[prefix] = append(groups[prefix], k)
	}

	sort.Strings(order)

	result := make([]string, 0, len(keys))
	for _, prefix := range order {
		group := groups[prefix]
		sort.Strings(group)
		result = append(result, group...)
	}
	return result
}

// extractPrefix returns the portion of a key before the first underscore,
// or the full key if no underscore is present.
func extractPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
