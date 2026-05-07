// Package rotator provides utilities for rotating keys in .env maps.
// It supports renaming keys to new names while optionally deprecating
// old keys by setting them to empty or removing them entirely.
package rotator

import "fmt"

// Strategy controls what happens to the old key after rotation.
type Strategy string

const (
	// StrategyRemove deletes the old key after copying its value.
	StrategyRemove Strategy = "remove"
	// StrategyDeprecate sets the old key to an empty string.
	StrategyDeprecate Strategy = "deprecate"
	// StrategyKeep leaves the old key unchanged.
	StrategyKeep Strategy = "keep"
)

// Options configures rotation behaviour.
type Options struct {
	// Strategy determines what happens to old keys.
	Strategy Strategy
	// FailOnMissing returns an error if an old key is not found.
	FailOnMissing bool
}

// DefaultOptions returns sensible rotation defaults.
func DefaultOptions() Options {
	return Options{
		Strategy:      StrategyRemove,
		FailOnMissing: false,
	}
}

// Rotation describes a single key rename operation.
type Rotation struct {
	OldKey string
	NewKey string
}

// Rotate applies a slice of Rotation rules to env, returning a new map.
func Rotate(env map[string]string, rotations []Rotation, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	for _, r := range rotations {
		val, exists := out[r.OldKey]
		if !exists {
			if opts.FailOnMissing {
				return nil, fmt.Errorf("rotator: key %q not found", r.OldKey)
			}
			continue
		}

		out[r.NewKey] = val

		switch opts.Strategy {
		case StrategyRemove:
			delete(out, r.OldKey)
		case StrategyDeprecate:
			out[r.OldKey] = ""
		case StrategyKeep:
			// leave old key as-is
		}
	}

	return out, nil
}
