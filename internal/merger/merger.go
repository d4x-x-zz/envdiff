// Package merger provides functionality to merge multiple .env files
// into a single unified map, with configurable conflict resolution strategies.
package merger

import "fmt"

// Strategy defines how key conflicts are resolved during merge.
type Strategy int

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last file that defines the key.
	StrategyLast
	// StrategyError returns an error if the same key appears with different values.
	StrategyError
)

// Options configures the merge behaviour.
type Options struct {
	Strategy Strategy
}

// DefaultOptions returns sensible defaults (first-wins).
func DefaultOptions() Options {
	return Options{Strategy: StrategyFirst}
}

// Merge combines multiple env maps into one according to the given options.
// The maps are applied in order, so maps[0] is considered the "first" file.
func Merge(maps []map[string]string, opts Options) (map[string]string, error) {
	result := make(map[string]string)

	for _, m := range maps {
		for k, v := range m {
			existing, exists := result[k]
			if !exists {
				result[k] = v
				continue
			}

			switch opts.Strategy {
			case StrategyFirst:
				// keep existing, do nothing
			case StrategyLast:
				result[k] = v
			case StrategyError:
				if existing != v {
					return nil, fmt.Errorf("merge conflict: key %q has values %q and %q", k, existing, v)
				}
			}
		}
	}

	return result, nil
}
