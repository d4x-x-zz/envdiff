// Package trimmer removes unused keys from a .env file by comparing
// it against a reference set of required keys.
package trimmer

import "sort"

// Options controls the behaviour of Trim.
type Options struct {
	// KeepUnknown retains keys that are not present in the required set
	// instead of removing them.
	KeepUnknown bool

	// Strict causes Trim to return an error when a required key is absent
	// from the source map.
	Strict bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		KeepUnknown: false,
		Strict:      false,
	}
}

// Result holds the output of a Trim operation.
type Result struct {
	// Trimmed is the filtered env map containing only retained keys.
	Trimmed map[string]string

	// Removed lists keys that were dropped from the source.
	Removed []string

	// Missing lists required keys that were absent from the source.
	Missing []string
}

// Trim filters src so that only keys present in required are kept.
// If opts.KeepUnknown is true, unknown keys are retained rather than removed.
// If opts.Strict is true and any required key is missing from src, an error
// is returned alongside a partial result.
func Trim(src map[string]string, required []string, opts Options) (Result, error) {
	reqSet := make(map[string]struct{}, len(required))
	for _, k := range required {
		reqSet[k] = struct{}{}
	}

	trimmed := make(map[string]string)
	var removed, missing []string

	for k, v := range src {
		if _, ok := reqSet[k]; ok {
			trimmed[k] = v
		} else if opts.KeepUnknown {
			trimmed[k] = v
		} else {
			removed = append(removed, k)
		}
	}

	for _, k := range required {
		if _, ok := src[k]; !ok {
			missing = append(missing, k)
		}
	}

	sort.Strings(removed)
	sort.Strings(missing)

	res := Result{Trimmed: trimmed, Removed: removed, Missing: missing}

	if opts.Strict && len(missing) > 0 {
		return res, fmt.Errorf("trimmer: %d required key(s) missing from source: %v", len(missing), missing)
	}

	return res, nil
}
