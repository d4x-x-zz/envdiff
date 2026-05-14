// Package caster infers and casts .env values to their likely Go types.
// It is useful for generating typed config structs or validating value formats.
package caster

import (
	"strconv"
	"strings"
)

// Type represents an inferred value type.
type Type string

const (
	TypeBool   Type = "bool"
	TypeInt    Type = "int"
	TypeFloat  Type = "float"
	TypeDSN    Type = "dsn"
	TypeURL    Type = "url"
	TypeString Type = "string"
)

// Result holds a key, its raw value, and the inferred type.
type Result struct {
	Key   string
	Value string
	Type  Type
}

// Options controls caster behaviour.
type Options struct {
	// DetectDSN enables heuristic detection of DSN-style values.
	DetectDSN bool
	// DetectURL enables heuristic detection of URL values.
	DetectURL bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		DetectDSN: true,
		DetectURL: true,
	}
}

// Cast inspects each value in env and returns a slice of typed Results.
func Cast(env map[string]string, opts Options) []Result {
	results := make([]Result, 0, len(env))
	for k, v := range env {
		results = append(results, Result{
			Key:   k,
			Value: v,
			Type:  infer(v, opts),
		})
	}
	return results
}

func infer(v string, opts Options) Type {
	if _, err := strconv.ParseBool(v); err == nil {
		return TypeBool
	}
	if _, err := strconv.ParseInt(v, 10, 64); err == nil {
		return TypeInt
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return TypeFloat
	}
	if opts.DetectURL {
		l := strings.ToLower(v)
		if strings.HasPrefix(l, "http://") || strings.HasPrefix(l, "https://") {
			return TypeURL
		}
	}
	if opts.DetectDSN {
		if isDSN(v) {
			return TypeDSN
		}
	}
	return TypeString
}

// isDSN uses a simple heuristic: scheme://user:pass@host/db
func isDSN(v string) bool {
	if idx := strings.Index(v, "://"); idx > 0 {
		rest := v[idx+3:]
		if strings.ContainsAny(rest, "@/") {
			return true
		}
	}
	return false
}
