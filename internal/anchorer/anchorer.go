// Package anchorer provides functionality to designate one .env file as the
// "anchor" and align all other files to its key set — flagging extras and
// filling in missing keys with a configurable default value.
package anchorer

import "sort"

// Options controls the behaviour of Anchor.
type Options struct {
	// FillMissing, when true, adds keys that exist in the anchor but are
	// absent in the target map, using DefaultValue as the value.
	FillMissing bool

	// RemoveExtra, when true, deletes keys from the target that do not exist
	// in the anchor.
	RemoveExtra bool

	// DefaultValue is used when FillMissing is true.
	DefaultValue string
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		FillMissing:  true,
		RemoveExtra:  false,
		DefaultValue: "",
	}
}

// Result holds the outcome of an Anchor call.
type Result struct {
	// Output is the (possibly modified) copy of the target map.
	Output map[string]string

	// Added lists keys that were inserted because they were missing.
	Added []string

	// Removed lists keys that were deleted because they were not in the anchor.
	Removed []string

	// Extra lists keys present in the target but absent in the anchor
	// (regardless of whether RemoveExtra was set).
	Extra []string
}

// Anchor aligns target against anchor according to opts.
func Anchor(anchorEnv, targetEnv map[string]string, opts Options) Result {
	output := make(map[string]string, len(targetEnv))
	for k, v := range targetEnv {
		output[k] = v
	}

	var added, removed, extra []string

	// Find keys missing from target.
	for k := range anchorEnv {
		if _, ok := output[k]; !ok {
			if opts.FillMissing {
				output[k] = opts.DefaultValue
				added = append(added, k)
			}
		}
	}

	// Find keys in target not present in anchor.
	for k := range targetEnv {
		if _, ok := anchorEnv[k]; !ok {
			extra = append(extra, k)
			if opts.RemoveExtra {
				delete(output, k)
				removed = append(removed, k)
			}
		}
	}

	sort.Strings(added)
	sort.Strings(removed)
	sort.Strings(extra)

	return Result{
		Output:  output,
		Added:   added,
		Removed: removed,
		Extra:   extra,
	}
}
