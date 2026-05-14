// Package splitter splits a flat env map into multiple maps based on key prefixes.
// Each resulting map contains only the keys belonging to that prefix group,
// with the prefix optionally stripped from the keys.
package splitter

import "sort"

// Options controls how splitting is performed.
type Options struct {
	// Prefixes is the list of prefixes to split on (e.g. ["DB_", "REDIS_"]).
	// Keys that don't match any prefix are collected under the "_other" group.
	Prefixes []string

	// StripPrefix removes the prefix from the key in the resulting map.
	StripPrefix bool

	// IncludeOther includes unmatched keys under the "_other" group.
	IncludeOther bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		StripPrefix:  true,
		IncludeOther: true,
	}
}

// Split partitions env into sub-maps keyed by prefix.
// Returns a map of prefix -> env map.
func Split(env map[string]string, opts Options) map[string]map[string]string {
	result := make(map[string]map[string]string)

	for k, v := range env {
		matched := false
		for _, prefix := range opts.Prefixes {
			if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
				if result[prefix] == nil {
					result[prefix] = make(map[string]string)
				}
				outKey := k
				if opts.StripPrefix {
					outKey = k[len(prefix):]
				}
				result[prefix][outKey] = v
				matched = true
				break
			}
		}
		if !matched && opts.IncludeOther {
			if result["_other"] == nil {
				result["_other"] = make(map[string]string)
			}
			result["_other"][k] = v
		}
	}
	return result
}

// SortedGroupNames returns the group names from a Split result in sorted order,
// with "_other" always last.
func SortedGroupNames(groups map[string]map[string]string) []string {
	names := make([]string, 0, len(groups))
	for k := range groups {
		if k != "_other" {
			names = append(names, k)
		}
	}
	sort.Strings(names)
	if _, ok := groups["_other"]; ok {
		names = append(names, "_other")
	}
	return names
}
