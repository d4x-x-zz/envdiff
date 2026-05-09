// Package interpolator resolves variable references within .env values.
// It expands expressions like ${VAR} or $VAR using values from the same
// map or a provided override map.
package interpolator

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var varPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Z_][A-Z0-9_]*)`)

// Options controls interpolation behaviour.
type Options struct {
	// FallbackToEnv allows falling back to OS environment variables
	// when a key is not found in the provided maps.
	FallbackToEnv bool
	// IgnoreMissing suppresses errors for unresolvable references.
	IgnoreMissing bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		FallbackToEnv: false,
		IgnoreMissing: true,
	}
}

// Interpolate resolves variable references in env values.
// The overrides map is consulted first, then the env map itself.
func Interpolate(env map[string]string, overrides map[string]string, opts Options) (map[string]string, error) {
	result := make(map[string]string, len(env))

	lookup := func(key string) (string, bool) {
		if overrides != nil {
			if v, ok := overrides[key]; ok {
				return v, true
			}
		}
		if v, ok := env[key]; ok {
			return v, true
		}
		if opts.FallbackToEnv {
			if v, ok := os.LookupEnv(key); ok {
				return v, true
			}
		}
		return "", false
	}

	for k, v := range env {
		expanded, err := expand(v, lookup, opts)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		result[k] = expanded
	}

	return result, nil
}

func expand(value string, lookup func(string) (string, bool), opts Options) (string, error) {
	var expandErr error
	result := varPattern.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		key := strings.TrimPrefix(strings.TrimPrefix(strings.Trim(match, "${}"), "${"), "$")
		key = strings.TrimSuffix(key, "}")
		if v, ok := lookup(key); ok {
			return v
		}
		if !opts.IgnoreMissing {
			expandErr = fmt.Errorf("unresolved variable: %s", key)
		}
		return match
	})
	return result, expandErr
}
