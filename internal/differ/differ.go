// Package differ compares two env maps and returns a structured result.
package differ

// Diff compares left and right env maps and returns a Result describing
// keys that are missing in either side or have mismatched values.
func Diff(left, right map[string]string) Result {
	var result Result

	// Check keys in left
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

	// Check keys only in right
	for k := range right {
		if _, ok := left[k]; !ok {
			result.MissingInLeft = append(result.MissingInLeft, k)
		}
	}

	return result
}
