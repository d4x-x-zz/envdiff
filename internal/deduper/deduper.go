// Package deduper removes duplicate key-value pairs across multiple env maps.
// When the same key appears with the same value in more than one source,
// the deduper can strip the redundant entries and report what was removed.
package deduper

// Options controls deduplication behaviour.
type Options struct {
	// KeepFirst retains the entry from the first map that defines a key.
	// When false the last definition wins (default: true).
	KeepFirst bool

	// SkipValueCheck treats two entries as duplicates only when the key
	// matches, regardless of value.
	SkipValueCheck bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{KeepFirst: true}
}

// Result holds the deduplicated map and a record of what was removed.
type Result struct {
	Env     map[string]string
	Removed []Duplicate
}

// Duplicate describes a single removed entry.
type Duplicate struct {
	Key        string
	Value      string
	SourceIndex int // index of the map the duplicate came from
}

// Dedupe merges the supplied maps and removes duplicate entries according to
// opts. Maps are processed in order; index 0 is the leftmost / primary source.
func Dedupe(maps []map[string]string, opts Options) Result {
	seen := make(map[string]string) // key → canonical value
	result := make(map[string]string)
	var removed []Duplicate

	for idx, m := range maps {
		for k, v := range m {
			canonical, exists := seen[k]
			if !exists {
				// First time we see this key — record it.
				seen[k] = v
				result[k] = v
				continue
			}

			// Key already seen — decide whether this is a duplicate.
			isDup := opts.SkipValueCheck || v == canonical
			if isDup {
				removed = append(removed, Duplicate{Key: k, Value: v, SourceIndex: idx})
				continue
			}

			// Different value — honour KeepFirst / KeepLast.
			if !opts.KeepFirst {
				seen[k] = v
				result[k] = v
			}
		}
	}

	return Result{Env: result, Removed: removed}
}
