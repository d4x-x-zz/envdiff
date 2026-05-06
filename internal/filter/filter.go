package filter

import "github.com/your-org/envdiff/internal/differ"

// Options controls which diff results are included.
type Options struct {
	// OnlyMissing restricts output to keys missing in either side.
	OnlyMissing bool
	// OnlyMismatched restricts output to keys present on both sides but with different values.
	OnlyMismatched bool
	// KeyPrefix, when non-empty, only includes keys that start with the given prefix.
	KeyPrefix string
}

// Apply returns a new Result containing only the entries that match opts.
func Apply(result differ.Result, opts Options) differ.Result {
	out := differ.Result{}

	for _, e := range result.MissingInRight {
		if opts.OnlyMismatched {
			continue
		}
		if opts.KeyPrefix != "" && !hasPrefix(e.Key, opts.KeyPrefix) {
			continue
		}
		out.MissingInRight = append(out.MissingInRight, e)
	}

	for _, e := range result.MissingInLeft {
		if opts.OnlyMismatched {
			continue
		}
		if opts.KeyPrefix != "" && !hasPrefix(e.Key, opts.KeyPrefix) {
			continue
		}
		out.MissingInLeft = append(out.MissingInLeft, e)
	}

	for _, e := range result.Mismatched {
		if opts.OnlyMissing {
			continue
		}
		if opts.KeyPrefix != "" && !hasPrefix(e.Key, opts.KeyPrefix) {
			continue
		}
		out.Mismatched = append(out.Mismatched, e)
	}

	return out
}

func hasPrefix(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}
