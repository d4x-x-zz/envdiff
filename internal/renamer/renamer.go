// Package renamer provides utilities for renaming keys across .env maps.
// It supports bulk renaming via a mapping of old->new key names, with options
// to control behaviour when a key is missing or a collision occurs.
package renamer

import "fmt"

// Strategy controls what happens when a rename target key already exists.
type Strategy int

const (
	// SkipOnConflict leaves the destination key untouched if it already exists.
	SkipOnConflict Strategy = iota
	// OverwriteOnConflict replaces the destination key with the renamed value.
	OverwriteOnConflict
	// ErrorOnConflict returns an error when the destination key already exists.
	ErrorOnConflict
)

// Options configures the Rename operation.
type Options struct {
	// Mapping is a map of oldKey -> newKey.
	Mapping map[string]string
	// IgnoreMissing skips keys in Mapping that don't exist in the source map.
	// When false, a missing source key returns an error.
	IgnoreMissing bool
	// ConflictStrategy determines behaviour on destination key collision.
	ConflictStrategy Strategy
}

// DefaultOptions returns an Options with safe defaults.
func DefaultOptions() Options {
	return Options{
		Mapping:          map[string]string{},
		IgnoreMissing:    true,
		ConflictStrategy: SkipOnConflict,
	}
}

// Rename applies the rename mapping to env, returning a new map with keys renamed.
// The original map is not modified.
func Rename(env map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	for oldKey, newKey := range opts.Mapping {
		val, exists := out[oldKey]
		if !exists {
			if opts.IgnoreMissing {
				continue
			}
			return nil, fmt.Errorf("renamer: source key %q not found", oldKey)
		}

		if _, conflict := out[newKey]; conflict && oldKey != newKey {
			switch opts.ConflictStrategy {
			case ErrorOnConflict:
				return nil, fmt.Errorf("renamer: destination key %q already exists", newKey)
			case SkipOnConflict:
				continue
			case OverwriteOnConflict:
				// fall through to assignment
			}
		}

		out[newKey] = val
		if oldKey != newKey {
			delete(out, oldKey)
		}
	}

	return out, nil
}
