// Package tagger assigns metadata tags to env keys based on pattern matching.
// Tags can be used downstream for filtering, reporting, or export grouping.
package tagger

import "strings"

// Options controls tagging behaviour.
type Options struct {
	// Rules maps a tag name to a list of key prefixes/suffixes that trigger it.
	Rules map[string][]string
	// CaseInsensitive makes prefix/suffix matching case-insensitive.
	CaseInsensitive bool
}

// DefaultOptions returns sensible defaults with common tag rules.
func DefaultOptions() Options {
	return Options{
		Rules: map[string][]string{
			"secret": {"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE"},
			"database": {"DB_", "DATABASE_", "POSTGRES", "MYSQL", "REDIS"},
			"network": {"HOST", "PORT", "URL", "ADDR", "ENDPOINT"},
			"feature": {"FEATURE_", "FLAG_", "ENABLE_", "DISABLE_"},
		},
		CaseInsensitive: true,
	}
}

// Result holds the tagging output for a single env map.
type Result struct {
	// Tags maps each key to its assigned tags (may be empty slice).
	Tags map[string][]string
}

// Tag assigns tags to every key in env according to opts.
func Tag(env map[string]string, opts Options) Result {
	result := Result{Tags: make(map[string][]string, len(env))}

	for key := range env {
		result.Tags[key] = matchTags(key, opts)
	}

	return result
}

func matchTags(key string, opts Options) []string {
	var tags []string
	cmp := key
	if opts.CaseInsensitive {
		cmp = strings.ToUpper(key)
	}

	for tag, patterns := range opts.Rules {
		for _, p := range patterns {
			pattern := p
			if opts.CaseInsensitive {
				pattern = strings.ToUpper(p)
			}
			if strings.Contains(cmp, pattern) {
				tags = append(tags, tag)
				break
			}
		}
	}

	return tags
}
