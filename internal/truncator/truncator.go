// Package truncator shortens env values that exceed a maximum length.
// Useful for logging or display contexts where long secrets or URLs
// should be capped to a readable size.
package truncator

import "strings"

// DefaultOptions returns a sensible set of truncation options.
func DefaultOptions() Options {
	return Options{
		MaxLen:   64,
		Suffix:   "...",
		KeysOnly: nil,
	}
}

// Options controls how truncation is applied.
type Options struct {
	// MaxLen is the maximum allowed value length (in runes) before truncation.
	MaxLen int

	// Suffix is appended to truncated values. Counts toward MaxLen.
	Suffix string

	// KeysOnly restricts truncation to the listed keys. If nil, all keys are
	// considered.
	KeysOnly []string
}

// Truncate applies length-capping to values in env according to opts.
// The original map is not modified; a new map is returned.
func Truncate(env map[string]string, opts Options) map[string]string {
	allowSet := make(map[string]struct{}, len(opts.KeysOnly))
	for _, k := range opts.KeysOnly {
		allowSet[strings.ToUpper(k)] = struct{}{}
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if len(allowSet) > 0 {
			if _, ok := allowSet[strings.ToUpper(k)]; !ok {
				out[k] = v
				continue
			}
		}
		out[k] = truncateValue(v, opts.MaxLen, opts.Suffix)
	}
	return out
}

func truncateValue(v string, maxLen int, suffix string) string {
	runes := []rune(v)
	if len(runes) <= maxLen {
		return v
	}
	suffixRunes := []rune(suffix)
	cutAt := maxLen - len(suffixRunes)
	if cutAt < 0 {
		cutAt = 0
	}
	return string(runes[:cutAt]) + suffix
}
