// Package redactor provides utilities for masking sensitive .env values
// based on key patterns or explicit lists.
package redactor

import "strings"

// Options controls redaction behaviour.
type Options struct {
	// SensitivePatterns is a list of substrings; any key containing one of
	// these (case-insensitive) will have its value redacted.
	SensitivePatterns []string
	// ExplicitKeys is an exact-match allowlist of keys to redact.
	ExplicitKeys []string
	// Mask is the string used to replace redacted values. Defaults to "***".
	Mask string
}

// DefaultOptions returns sensible defaults that catch common secrets.
func DefaultOptions() Options {
	return Options{
		SensitivePatterns: []string{"secret", "password", "passwd", "token", "key", "auth", "credential", "private"},
		Mask:              "***",
	}
}

// Redact returns a copy of env with sensitive values replaced by the mask.
func Redact(env map[string]string, opts Options) map[string]string {
	if opts.Mask == "" {
		opts.Mask = "***"
	}

	explicit := make(map[string]bool, len(opts.ExplicitKeys))
	for _, k := range opts.ExplicitKeys {
		explicit[k] = true
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if explicit[k] || isSensitive(k, opts.SensitivePatterns) {
			out[k] = opts.Mask
		} else {
			out[k] = v
		}
	}
	return out
}

// isSensitive returns true when key contains any of the given patterns
// (case-insensitive).
func isSensitive(key string, patterns []string) bool {
	lower := strings.ToLower(key)
	for _, p := range patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}
