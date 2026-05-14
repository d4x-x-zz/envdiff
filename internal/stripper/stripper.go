// Package stripper removes keys from an env map based on patterns or an explicit list.
package stripper

import "strings"

// Options controls which keys are stripped.
type Options struct {
	// Keys is an explicit list of key names to remove.
	Keys []string
	// Prefixes removes any key whose name starts with one of these prefixes.
	Prefixes []string
	// Suffixes removes any key whose name ends with one of these suffixes.
	Suffixes []string
	// DryRun returns the list of keys that would be removed without mutating the map.
	DryRun bool
}

// DefaultOptions returns an Options with no rules set.
func DefaultOptions() Options {
	return Options{}
}

// Strip removes matching keys from env and returns the modified map plus the
// list of keys that were (or would be) removed.
func Strip(env map[string]string, opts Options) (map[string]string, []string) {
	removed := []string{}

	should := func(key string) bool {
		for _, k := range opts.Keys {
			if key == k {
				return true
			}
		}
		for _, p := range opts.Prefixes {
			if strings.HasPrefix(key, p) {
				return true
			}
		}
		for _, s := range opts.Suffixes {
			if strings.HasSuffix(key, s) {
				return true
			}
		}
		return false
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if should(k) {
			removed = append(removed, k)
		} else {
			out[k] = v
		}
	}

	if opts.DryRun {
		return env, removed
	}
	return out, removed
}
