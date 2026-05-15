// Package extractor pulls a subset of keys from an env map based on
// an explicit allow-list or a glob-style pattern.
package extractor

import (
	"path"
	"sort"
)

// Options controls how Extract behaves.
type Options struct {
	// Keys is an explicit list of keys to extract. If non-empty, Patterns is
	// ignored.
	Keys []string

	// Patterns is a list of glob patterns (e.g. "DB_*") matched against key
	// names. Used only when Keys is empty.
	Patterns []string

	// IgnoreMissing silently skips keys that do not exist in src when true.
	// When false, missing explicit keys are still skipped (no error is
	// returned), but the returned Missed slice will be populated.
	IgnoreMissing bool
}

// DefaultOptions returns a safe default configuration.
func DefaultOptions() Options {
	return Options{
		IgnoreMissing: true,
	}
}

// Result holds the output of an Extract call.
type Result struct {
	// Env contains the extracted key-value pairs.
	Env map[string]string
	// Missed lists keys that were requested but not found in src.
	Missed []string
}

// Extract returns a new map containing only the keys selected by opts.
func Extract(src map[string]string, opts Options) Result {
	out := make(map[string]string)
	var missed []string

	if len(opts.Keys) > 0 {
		for _, k := range opts.Keys {
			v, ok := src[k]
			if !ok {
				missed = append(missed, k)
				continue
			}
			out[k] = v
		}
	} else {
		for k, v := range src {
			if matchesAny(k, opts.Patterns) {
				out[k] = v
			}
		}
	}

	sort.Strings(missed)
	return Result{Env: out, Missed: missed}
}

func matchesAny(key string, patterns []string) bool {
	for _, p := range patterns {
		if ok, _ := path.Match(p, key); ok {
			return true
		}
	}
	return false
}
