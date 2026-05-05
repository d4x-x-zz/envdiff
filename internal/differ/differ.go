package differ

// Result holds the diff outcome between two env files.
type Result struct {
	// MissingInRight are keys present in left but absent in right.
	MissingInRight []string
	// MissingInLeft are keys present in right but absent in left.
	MissingInLeft []string
	// Mismatched are keys present in both files but with different values.
	Mismatched []MismatchedKey
}

// MismatchedKey captures a key whose value differs between the two files.
type MismatchedKey struct {
	Key        string
	LeftValue  string
	RightValue string
}

// HasDiff returns true when any difference was found.
func (r Result) HasDiff() bool {
	return len(r.MissingInRight) > 0 ||
		len(r.MissingInLeft) > 0 ||
		len(r.Mismatched) > 0
}

// Diff compares two parsed env maps and returns a Result describing
// all differences. The maps are typically produced by parser.ParseFile.
func Diff(left, right map[string]string) Result {
	var result Result

	// Keys in left — check presence and value equality in right.
	for k, lv := range left {
		rv, ok := right[k]
		if !ok {
			result.MissingInRight = append(result.MissingInRight, k)
			continue
		}
		if lv != rv {
			result.Mismatched = append(result.Mismatched, MismatchedKey{
				Key:        k,
				LeftValue:  lv,
				RightValue: rv,
			})
		}
	}

	// Keys in right that are absent in left.
	for k := range right {
		if _, ok := left[k]; !ok {
			result.MissingInLeft = append(result.MissingInLeft, k)
		}
	}

	return result
}
