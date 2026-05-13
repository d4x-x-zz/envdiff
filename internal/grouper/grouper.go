// Package grouper organises a flat env map into named groups based on key
// prefix conventions (e.g. DB_HOST → group "DB").
package grouper

import (
	"sort"
	"strings"
)

// DefaultOptions returns a GroupOptions with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Separator:    "_",
		UngroupedKey: "OTHER",
	}
}

// Options controls how grouping is performed.
type Options struct {
	// Separator is the delimiter used to split a key into prefix + rest.
	Separator string
	// UngroupedKey is the group name assigned to keys with no prefix.
	UngroupedKey string
	// AllowList, when non-empty, limits grouping to these prefixes only.
	// Keys whose prefix is not in AllowList fall into UngroupedKey.
	AllowList []string
}

// Group partitions env into a map of prefix → (key → value).
// Keys that have no separator are placed under Options.UngroupedKey.
func Group(env map[string]string, opts Options) map[string]map[string]string {
	allowed := make(map[string]bool, len(opts.AllowList))
	for _, p := range opts.AllowList {
		allowed[strings.ToUpper(p)] = true
	}

	result := make(map[string]map[string]string)

	for k, v := range env {
		prefix := extractPrefix(k, opts.Separator)
		if prefix == "" {
			prefix = opts.UngroupedKey
		} else if len(allowed) > 0 && !allowed[prefix] {
			prefix = opts.UngroupedKey
		}
		if result[prefix] == nil {
			result[prefix] = make(map[string]string)
		}
		result[prefix][k] = v
	}
	return result
}

// SortedGroupNames returns group names in alphabetical order.
func SortedGroupNames(groups map[string]map[string]string) []string {
	names := make([]string, 0, len(groups))
	for n := range groups {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func extractPrefix(key, sep string) string {
	idx := strings.Index(key, sep)
	if idx <= 0 {
		return ""
	}
	return key[:idx]
}
