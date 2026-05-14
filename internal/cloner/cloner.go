// Package cloner copies an env map and optionally transforms keys or values
// during the copy operation.
package cloner

import "strings"

// Options controls how the clone is performed.
type Options struct {
	// KeyPrefix adds a prefix to every key in the cloned map.
	KeyPrefix string
	// KeySuffix adds a suffix to every key in the cloned map.
	KeySuffix string
	// UppercaseKeys converts all keys to uppercase.
	UppercaseKeys bool
	// OmitEmpty skips keys whose value is an empty string.
	OmitEmpty bool
	// ValueTransform is an optional function applied to every value.
	ValueTransform func(string) string
}

// DefaultOptions returns an Options with no transformations applied.
func DefaultOptions() Options {
	return Options{}
}

// Clone returns a deep copy of src, applying any transformations defined in opts.
func Clone(src map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		if opts.OmitEmpty && v == "" {
			continue
		}

		key := k
		if opts.UppercaseKeys {
			key = strings.ToUpper(key)
		}
		key = opts.KeyPrefix + key + opts.KeySuffix

		val := v
		if opts.ValueTransform != nil {
			val = opts.ValueTransform(val)
		}

		out[key] = val
	}
	return out
}
