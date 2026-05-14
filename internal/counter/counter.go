// Package counter provides key counting and frequency analysis for env maps.
package counter

import (
	"sort"
	"strings"
)

// Options configures the Count operation.
type Options struct {
	// KeyPrefix filters keys to only those starting with the given prefix.
	KeyPrefix string
	// CaseSensitive controls whether key matching is case-sensitive.
	CaseSensitive bool
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		CaseSensitive: true,
	}
}

// Entry holds a key and the number of times its value appears across all maps.
type Entry struct {
	Key   string
	Count int
}

// Result holds the output of a Count operation.
type Result struct {
	// Total is the number of distinct keys found.
	Total int
	// Entries lists each key and how many of the provided maps contain it.
	Entries []Entry
}

// Count tallies how many of the provided env maps contain each key.
// Maps with zero entries are ignored. Keys are returned sorted alphabetically.
func Count(maps []map[string]string, opts Options) Result {
	freq := make(map[string]int)

	for _, m := range maps {
		for k := range m {
			key := k
			if !opts.CaseSensitive {
				key = strings.ToUpper(k)
			}
			if opts.KeyPrefix != "" && !strings.HasPrefix(key, opts.KeyPrefix) {
				continue
			}
			freq[key]++
		}
	}

	entries := make([]Entry, 0, len(freq))
	for k, c := range freq {
		entries = append(entries, Entry{Key: k, Count: c})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return Result{
		Total:   len(entries),
		Entries: entries,
	}
}
