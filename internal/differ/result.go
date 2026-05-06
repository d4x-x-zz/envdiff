package differ

// Result holds the outcome of diffing two env maps.
type Result struct {
	// MissingInRight contains keys present in left but absent in right.
	MissingInRight []string
	// MissingInLeft contains keys present in right but absent in left.
	MissingInLeft []string
	// Mismatched contains keys present in both files but with different values.
	Mismatched []MismatchEntry
}

// MismatchEntry records a key whose value differs between files.
type MismatchEntry struct {
	Key        string
	LeftValue  string
	RightValue string
}

// Clean returns true when there are no differences of any kind.
func (r Result) Clean() bool {
	return len(r.MissingInRight) == 0 &&
		len(r.MissingInLeft) == 0 &&
		len(r.Mismatched) == 0
}

// TotalIssues returns the total count of all detected differences.
func (r Result) TotalIssues() int {
	return len(r.MissingInRight) + len(r.MissingInLeft) + len(r.Mismatched)
}

// AllKeys returns a deduplicated, sorted list of all keys involved in any difference.
func (r Result) AllKeys() []string {
	seen := make(map[string]struct{})
	for _, k := range r.MissingInRight {
		seen[k] = struct{}{}
	}
	for _, k := range r.MissingInLeft {
		seen[k] = struct{}{}
	}
	for _, e := range r.Mismatched {
		seen[e.Key] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
