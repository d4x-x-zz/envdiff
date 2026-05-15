// Package aliaser provides functionality to create key aliases in an env map.
// An alias copies the value of an existing key to one or more new keys.
package aliaser

import "fmt"

// Options controls aliaser behaviour.
type Options struct {
	// Aliases maps source keys to one or more destination keys.
	Aliases map[string][]string
	// OverwriteExisting allows destination keys that already exist to be overwritten.
	OverwriteExisting bool
	// IgnoreMissing skips source keys that are not present instead of returning an error.
	IgnoreMissing bool
}

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Aliases:           map[string][]string{},
		OverwriteExisting: false,
		IgnoreMissing:     true,
	}
}

// Alias applies the alias rules in opts to env, returning a new map with the
// aliased keys added. The original map is not modified.
func Alias(env map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	for src, dests := range opts.Aliases {
		val, ok := env[src]
		if !ok {
			if opts.IgnoreMissing {
				continue
			}
			return nil, fmt.Errorf("aliaser: source key %q not found", src)
		}
		for _, dest := range dests {
			if _, exists := out[dest]; exists && !opts.OverwriteExisting {
				continue
			}
			out[dest] = val
		}
	}

	return out, nil
}
