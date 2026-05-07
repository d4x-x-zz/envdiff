// Package patcher applies a set of key-value patches to an existing env map.
// It supports adding new keys, updating existing ones, and deleting keys.
package patcher

import "fmt"

// Op represents a patch operation type.
type Op string

const (
	OpSet    Op = "set"
	OpDelete Op = "delete"
)

// Patch describes a single change to apply to an env map.
type Patch struct {
	Op    Op
	Key   string
	Value string // only used for OpSet
}

// Options controls patcher behaviour.
type Options struct {
	// ErrorOnMissingDelete causes Patch to return an error when deleting a key
	// that does not exist in the source map.
	ErrorOnMissingDelete bool
	// ErrorOnNoChange causes Patch to return an error when a set operation
	// would not change the current value.
	ErrorOnNoChange bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		ErrorOnMissingDelete: false,
		ErrorOnNoChange:      false,
	}
}

// Apply applies patches to a copy of src and returns the resulting map.
func Apply(src map[string]string, patches []Patch, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}

	for _, p := range patches {
		switch p.Op {
		case OpSet:
			if opts.ErrorOnNoChange {
				if cur, exists := out[p.Key]; exists && cur == p.Value {
					return nil, fmt.Errorf("patch: key %q already has value %q", p.Key, p.Value)
				}
			}
			out[p.Key] = p.Value
		case OpDelete:
			if _, exists := out[p.Key]; !exists {
				if opts.ErrorOnMissingDelete {
					return nil, fmt.Errorf("patch: cannot delete missing key %q", p.Key)
				}
				continue
			}
			delete(out, p.Key)
		default:
			return nil, fmt.Errorf("patch: unknown op %q for key %q", p.Op, p.Key)
		}
	}

	return out, nil
}
