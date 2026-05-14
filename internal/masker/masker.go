// Package masker provides utilities for partially masking env values,
// preserving a configurable number of characters at the start or end.
package masker

import "strings"

// Options controls masking behaviour.
type Options struct {
	// MaskChar is the character used to replace hidden characters. Default: "*".
	MaskChar string
	// VisiblePrefix is the number of leading characters to keep visible. Default: 2.
	VisiblePrefix int
	// VisibleSuffix is the number of trailing characters to keep visible. Default: 0.
	VisibleSuffix int
	// MinLength is the minimum value length required before masking is applied.
	// Shorter values are fully masked. Default: 4.
	MinLength int
}

// DefaultOptions returns sensible masking defaults.
func DefaultOptions() Options {
	return Options{
		MaskChar:      "*",
		VisiblePrefix: 2,
		VisibleSuffix: 0,
		MinLength:     4,
	}
}

// Mask applies partial masking to every value in env according to opts.
// Keys are never modified.
func Mask(env map[string]string, opts Options) map[string]string {
	if opts.MaskChar == "" {
		opts.MaskChar = "*"
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = maskValue(v, opts)
	}
	return out
}

func maskValue(v string, opts Options) string {
	if len(v) < opts.MinLength {
		return strings.Repeat(opts.MaskChar, len(v))
	}

	prefix := opts.VisiblePrefix
	if prefix > len(v) {
		prefix = len(v)
	}

	suffix := opts.VisibleSuffix
	if prefix+suffix >= len(v) {
		suffix = 0
	}

	hidden := len(v) - prefix - suffix
	if hidden <= 0 {
		return v
	}

	var sb strings.Builder
	sb.WriteString(v[:prefix])
	sb.WriteString(strings.Repeat(opts.MaskChar, hidden))
	if suffix > 0 {
		sb.WriteString(v[len(v)-suffix:])
	}
	return sb.String()
}
