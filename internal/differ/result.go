package differ

// Side indicates which file a key is missing from.
type Side string

const (
	SideLeft  Side = "left"
	SideRight Side = "right"
)

// Diff represents a single key difference between two env files.
type Diff struct {
	Key      string
	Side     Side   // populated for missing keys
	LeftVal  string // populated for mismatched keys
	RightVal string // populated for mismatched keys
}

// Result holds the full comparison output between two env files.
type Result struct {
	Missing    []Diff
	Mismatched []Diff
}

// Clean returns true when there are no differences.
func (r *Result) Clean() bool {
	return len(r.Missing) == 0 && len(r.Mismatched) == 0
}

// TotalIssues returns the combined count of missing and mismatched keys.
func (r *Result) TotalIssues() int {
	return len(r.Missing) + len(r.Mismatched)
}

// MissingKeys returns only the keys that are absent from one side.
func (r *Result) MissingKeys() []string {
	keys := make([]string, 0, len(r.Missing))
	for _, d := range r.Missing {
		keys = append(keys, d.Key)
	}
	return keys
}

// MismatchedKeys returns only the keys whose values differ.
func (r *Result) MismatchedKeys() []string {
	keys := make([]string, 0, len(r.Mismatched))
	for _, d := range r.Mismatched {
		keys = append(keys, d.Key)
	}
	return keys
}
