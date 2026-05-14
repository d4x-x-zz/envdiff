// Package flattener collapses nested key structures (e.g. APP__DB__HOST)
// into a single-level map using a configurable separator.
package flattener

import "strings"

// Options controls how flattening behaves.
type Options struct {
	// Separator is the delimiter used to detect nesting. Defaults to "__".
	Separator string
	// Depth is the maximum number of segments to split on. 0 means unlimited.
	Depth int
	// LowercaseKeys converts all resulting keys to lowercase.
	LowercaseKeys bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Separator:     "__",
		Depth:         0,
		LowercaseKeys: false,
	}
}

// Flatten takes an env map and returns a new map where each key is split
// on the separator and reassembled with a single dot, up to Depth levels.
// Keys that do not contain the separator are passed through unchanged.
func Flatten(env map[string]string, opts Options) map[string]string {
	if opts.Separator == "" {
		opts.Separator = "__"
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		newKey := flattenKey(k, opts)
		out[newKey] = v
	}
	return out
}

func flattenKey(key string, opts Options) string {
	parts := strings.Split(key, opts.Separator)

	if opts.Depth > 0 && len(parts) > opts.Depth+1 {
		// rejoin the tail beyond Depth back with the separator
		head := parts[:opts.Depth]
		tail := strings.Join(parts[opts.Depth:], opts.Separator)
		parts = append(head, tail)
	}

	result := strings.Join(parts, ".")
	if opts.LowercaseKeys {
		result = strings.ToLower(result)
	}
	return result
}
